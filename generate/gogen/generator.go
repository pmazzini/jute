package gogen

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-zookeeper/jute/generate"
	"github.com/go-zookeeper/jute/parser"
)

// ModuleMap helps map jute modules to go packages.
type ModuleMap struct {
	Re   *regexp.Regexp // Regexp to match
	Repl string         // Re replacement string (see regexp.ReplaceAll)
	Skip bool           // Don't generate files for the matching jute modules
}

func (m ModuleMap) String() string {
	sb := strings.Builder{}
	sb.WriteString(m.Re.String())
	sb.WriteString(":")

	if m.Skip {
		sb.WriteString("-")
	} else {
		sb.WriteString(m.Repl)
	}

	return sb.String()

}

// Options can be passed to Generate to modify the output.
type Options struct {
	// Base import path for generated packages
	ImportPathPrefix string
	// if not blank will use for serialization/deserialization libraries
	JuteImport string
	// map of hute modles to go packages. The BaseImportPath is prepended to
	// these values.
	ModuleMap []ModuleMap
}

func defaultOptions() *Options {
	return &Options{
		JuteImport: "github.com/go-zookeeper/jute/lib/go/jute",
	}
}

// Generate will generate Go packages/modules for the given input files.
func Generate(outDir string, files []*generate.File, opts *Options) error {
	if opts == nil {
		opts = defaultOptions()
	}
	g := &generator{
		opts:      opts,
		outDir:    outDir,
		moduleMap: make(map[string]goPackage),
	}

	for _, file := range files {
		for _, m := range file.Doc.Modules {
			g.addModule(file.Path, m)
		}
	}

	return g.generate()
}

type module struct {
	node *parser.Module // parsed module

	srcFilename string // abs path to source jute filename
	goPkg       goPackage
	classes     []*class
}

type class struct {
	node *parser.Class

	goName        string
	importModules []string // jute modules required for import
	fields        []*field
}

func (cls *class) hasContainers() bool {
	for _, f := range cls.fields {
		switch f.node.Type.(type) {
		case *parser.MapType, *parser.VectorType:
			return true
		}
	}
	return false
}

type field struct {
	node   *parser.Field
	goName string
	goType string
}

type generator struct {
	opts    *Options
	outDir  string
	modules []*module

	moduleMap map[string]goPackage // map of jute module name to go import path
}

// Add module will add a module to the generator adding some
// specific go metadata along the way.
func (g *generator) addModule(srcFilename string, node *parser.Module) error {
	goPkg, ok := g.goPkg(node.Name)
	if !ok {
		log.Printf("skipping module %s", node.Name)
		return nil
	}

	g.moduleMap[node.Name] = goPkg
	m := &module{
		node:        node,
		srcFilename: srcFilename,
		goPkg:       goPkg,
	}

	for _, classNode := range node.Classes {
		cls := &class{
			node:   classNode,
			goName: camelcase(classNode.Name),
		}

		for _, fieldNode := range classNode.Fields {
			typ, err := g.goType(fieldNode.Type)
			if err != nil {
				return err
			}

			if t, ok := fieldNode.Type.(*parser.ClassType); ok && t.Namespace != "" {
				cls.importModules = append(cls.importModules, t.Namespace)
			}

			fld := &field{
				node:   fieldNode,
				goName: camelcase(fieldNode.Name),
				goType: typ,
			}
			cls.fields = append(cls.fields, fld)
		}
		m.classes = append(m.classes, cls)
	}

	//	g.moduleMap[node.Name] = m.goImportPath
	g.modules = append(g.modules, m)
	return nil
}

func (g *generator) generate() error {
	for _, m := range g.modules {
		if err := g.writeModule(m); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) writeModule(m *module) error {
	dir := filepath.Join(g.outDir, m.goPkg.relPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create dir '%s': %w", dir, err)
	}

	for _, cls := range m.classes {
		filename := filepath.Join(dir, strings.ToLower(cls.node.Name)+".go")
		log.Printf("writing %s", filename)
		fw := &fileWriter{}

		g.writeHeader(fw, m.srcFilename, m.goPkg.name, m.goPkg.importPath)

		imports := []string{}
		for _, imp := range cls.importModules {
			pkg := g.moduleMap[imp]
			imports = append(imports, pkg.importPath)
		}

		g.writeImports(fw, imports)
		g.writeClassStruct(fw, cls)
		g.writeReadMethod(fw, cls)
		g.writeWriteMethod(fw, cls)
		g.writeStringMethod(fw, cls)

		err := fw.writeFile(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) writeHeader(fw *fileWriter, srcFilename, pkg, importPath string) {
	fw.printf("// Autogenerated jute compiler\n")
	fw.printf("// @generated from '%s'\n", srcFilename)
	fw.printf("\n")
	fw.printf("package %s // %s\n", pkg, importPath)
	fw.printf("\n")
}

func (g *generator) writeImports(fw *fileWriter, imports []string) {
	fw.printf("import (\n")
	fw.printf("\t\"fmt\"\n")
	fw.printf("\n")
	fw.printf("\t\"%s\"\n", g.opts.JuteImport)
	for _, imp := range imports {
		fw.printf("\t\"%s\"\n", imp)
	}
	fw.printf(")\n\n")
}

func (g *generator) writeClassStruct(fw *fileWriter, cls *class) {
	fw.printf("type %s struct {\n", cls.goName)
	for _, fld := range cls.fields {
		fw.printf("\t%s ", fld.goName)
		fw.printf("%s // %s\n", fld.goType, fld.node.Name)
	}
	fw.printf("}\n\n")
}

func (g *generator) writeWriteMethod(fw *fileWriter, cls *class) {
	fw.printf("func (r *%s) Write(enc jute.Encoder) error {\n", cls.goName)
	fw.printf("\tif err := enc.WriteStart(); err != nil {\n")
	fw.printf("\t\treturn err\n")
	fw.printf("\t}\n")

	for _, fld := range cls.fields {
		method, err := g.serializeMethod(fld.node.Type, "r."+fld.goName)
		if err != nil {
			panic(err)
		}
		fw.printf(method)
	}

	fw.printf("\tif err := enc.WriteEnd(); err != nil {\n")
	fw.printf("\t\treturn err\n")
	fw.printf("\t}\n")
	fw.printf("\treturn nil\n")
	fw.printf("}\n\n")
}

func (g *generator) writeReadMethod(fw *fileWriter, cls *class) {
	fw.printf("func (r *%s) Read(dec jute.Decoder) (err error) {\n", cls.goName)
	// create a size variable if we have maps or vectors
	if cls.hasContainers() {
		fw.printf("\tvar size int\n")
	}
	fw.printf("\tif err = dec.ReadStart(); err != nil {\n")
	fw.printf("\t\treturn err\n")
	fw.printf("\t}\n")
	for i, fld := range cls.fields {
		method, err := g.deserializeMethod(fld.node.Type, "r."+fld.goName, i)
		if err != nil {
			panic(err)
		}
		fw.printf(method)
	}
	fw.printf("\tif err = dec.ReadEnd(); err != nil {\n")
	fw.printf("\t\treturn err\n")
	fw.printf("\t}\n")
	fw.printf("\treturn nil\n")
	fw.printf("}\n\n")
}

func (g *generator) writeStringMethod(fw *fileWriter, cls *class) {
	fw.printf("func (r *%s) String() string {\n", cls.goName)
	fw.printf("\tif r == nil {\n")
	fw.printf("\t\treturn \"<nil>\"\n")
	fw.printf("\t}\n")
	fw.printf("return fmt.Sprintf(\"%s(%%+v)\", *r)\n", cls.goName)
	fw.printf("\t}\n\n")
}

func (g *generator) serializeMethod(juteType parser.Type, fieldName string) (string, error) {
	w := &strings.Builder{}
	switch t := juteType.(type) {
	case *parser.PType:
		typeName, ok := primTypeName[t.TypeID]
		if !ok {
			return "", fmt.Errorf("unknown primative type: %v", t.TypeID)
		}
		fmt.Fprintf(w, "if err := enc.Write%s(%s); err != nil {\n", typeName, fieldName)
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
	case *parser.VectorType:
		itemMethod, err := g.serializeMethod(t.Type, "v")
		if err != nil {
			return "", err
		}

		fmt.Fprintf(w, "if err := enc.WriteVectorStart(len(%s)); err != nil {\n", fieldName)
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "for _, v := range %s {\n", fieldName)
		fmt.Fprintf(w, itemMethod)
		fmt.Fprintf(w, "}\n") // end for loop
		fmt.Fprintf(w, "if err := enc.WriteVectorEnd(); err != nil {\n")
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
	case *parser.MapType:
		keyMethod, err := g.serializeMethod(t.KeyType, "k")
		if err != nil {
			return "", err
		}

		valMethod, err := g.serializeMethod(t.ValType, "v")
		if err != nil {
			return "", err
		}

		fmt.Fprintf(w, "if err := enc.WriteMapStart(len(%s)); err != nil {\n", fieldName)
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "for k, v := range %s {\n", fieldName)
		fmt.Fprintf(w, keyMethod)
		fmt.Fprintf(w, valMethod)
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "if err := enc.WriteMapEnd(); err != nil {\n")
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
	case *parser.ClassType:
		fmt.Fprintf(w, "if err := enc.WriteRecord(%s); err != nil {\n", fieldName)
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
	default:
		return "", fmt.Errorf("unknown type %T for field %s", t, fieldName)
	}
	return w.String(), nil
}

func (g *generator) deserializeMethod(juteType parser.Type, fieldName string, idx int) (string, error) {
	w := &strings.Builder{}
	switch t := juteType.(type) {
	case *parser.PType:
		typeName, ok := primTypeName[t.TypeID]
		if !ok {
			return "", fmt.Errorf("unknown primative type: %v", t.TypeID)
		}
		fmt.Fprintf(w, "%s, err = dec.Read%s()\n", fieldName, typeName)
		fmt.Fprintf(w, "if err != nil {\n")
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
	case *parser.VectorType:
		itemType, err := g.goType(juteType)
		if err != nil {
			return "", err
		}
		itemMethod, err := g.deserializeMethod(t.Type, fmt.Sprintf("%s[i]", fieldName), idx+1)
		if err != nil {
			return "", err
		}

		fmt.Fprintf(w, "size, err = dec.ReadVectorStart()\n")
		fmt.Fprintf(w, "if err != nil {\n")
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "\t%s = make(%s, size)\n", fieldName, itemType)
		fmt.Fprintf(w, "\tfor i := 0; i < size; i++ {\n")
		fmt.Fprintf(w, itemMethod)
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "if err = dec.ReadVectorEnd(); err != nil {\n")
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
	case *parser.MapType:
		mapType, err := g.goType(juteType)
		if err != nil {
			return "", err
		}

		keytype, err := g.goType(t.KeyType)
		if err != nil {
			return "", err
		}

		valtype, err := g.goType(t.ValType)
		if err != nil {
			return "", err
		}

		keyMethod, err := g.deserializeMethod(t.KeyType, fmt.Sprintf("k%d", idx), idx+1)
		if err != nil {
			return "", err
		}

		valMethod, err := g.deserializeMethod(t.ValType, fmt.Sprintf("v%d", idx), idx+1)
		if err != nil {
			return "", err
		}

		fmt.Fprintf(w, "size, err = dec.ReadMapStart()\n")
		fmt.Fprintf(w, "if err != nil {\n")
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "%s = make(%s)\n", fieldName, mapType)
		fmt.Fprintf(w, "var k%d %s\n", idx, keytype)
		fmt.Fprintf(w, "var v%d %s\n", idx, valtype)
		fmt.Fprintf(w, "for i := 0; i < size; i++ {\n")
		fmt.Fprintf(w, keyMethod)
		fmt.Fprintf(w, valMethod)
		fmt.Fprintf(w, "\t%s[k%d] = v%d\n", fieldName, idx, idx)
		fmt.Fprintf(w, "}\n")
		fmt.Fprintf(w, "if err = dec.ReadMapEnd(); err != nil {\n")
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
	case *parser.ClassType:
		fmt.Fprintf(w, "if err = dec.ReadRecord(%s); err != nil {\n", fieldName)
		fmt.Fprintf(w, "\treturn err\n")
		fmt.Fprintf(w, "}\n")
	default:
		return "", fmt.Errorf("unknown type %T for field %s", t, fieldName)
	}
	return w.String(), nil
}

// goType will go type as string for the given jute ast type.
func (g *generator) goType(juteType parser.Type) (string, error) {
	switch t := juteType.(type) {
	case *parser.PType:
		if goType, ok := primaryTypeMap[t.TypeID]; ok {
			return goType, nil
		}
		return "", fmt.Errorf("unknown primative type %v", t.TypeID)

	case *parser.VectorType:
		innerType, err := g.goType(t.Type)
		if err != nil {
			return "", err
		}
		return "[]" + innerType, nil

	case *parser.MapType:
		keyType, err := g.goType(t.KeyType)
		if err != nil {
			return "", err
		}

		valType, err := g.goType(t.ValType)
		if err != nil {
			return "", err
		}

		return "map[" + keyType + "]" + valType, nil
	case *parser.ClassType:
		var pkg string
		if t.Namespace != "" {
			pkg = g.moduleMap[t.Namespace].name + "."
		}
		return "*" + pkg + t.ClassName, nil
	}
	return "", fmt.Errorf("unknown type %T", juteType)
}

func camelcase(ident string) string {
	var out string
	var upper bool
	for i, r := range ident {
		if i == 0 {
			out += string(unicode.ToUpper(r))
			continue
		}

		if upper {
			out += string(unicode.ToUpper(r))
			upper = false
			continue
		}

		if r == '_' {
			upper = true
			continue
		}

		out += string(r)
	}

	return out
}

// pkgName will return a package name and a import path from a jute module
// name.
func (g *generator) goPkg(module string) (goPackage, bool) {
	importPath := module
	for _, mm := range g.opts.ModuleMap {
		if mm.Skip && mm.Re.MatchString(module) {
			return goPackage{}, false
		}
		importPath = mm.Re.ReplaceAllString(importPath, mm.Repl)
	}

	// clean up any hanging slashes/dots
	importPath = strings.Trim(importPath, "/.")
	importPath = strings.ReplaceAll(importPath, ".", "/")

	var pkgName string
	i := strings.LastIndexAny(importPath, "./")
	if i > 0 {
		pkgName = importPath[i+1:]
	} else {
		pkgName = importPath
	}

	return goPackage{
		name:       pkgName,
		importPath: path.Join(g.opts.ImportPathPrefix, importPath),
		relPath:    importPath,
	}, true
}

type goPackage struct {
	name       string
	importPath string
	relPath    string
}
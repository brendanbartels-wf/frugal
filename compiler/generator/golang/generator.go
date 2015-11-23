package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Workiva/frugal/compiler/generator"
	"github.com/Workiva/frugal/compiler/globals"
	"github.com/Workiva/frugal/compiler/parser"
)

const (
	lang             = "go"
	defaultOutputDir = "gen-go"
)

type Generator struct {
	*generator.BaseGenerator
	generateConstants bool
}

func NewGenerator(options map[string]string) generator.MultipleFileGenerator {
	return &Generator{&generator.BaseGenerator{Options: options}, true}
}

func (g *Generator) GetOutputDir(dir string, f *parser.Frugal) string {
	if pkg, ok := f.Thrift.Namespaces[lang]; ok {
		path := generator.GetPackageComponents(pkg)
		dir = filepath.Join(append([]string{dir}, path...)...)
	} else {
		dir = filepath.Join(dir, f.Name)
	}
	return dir
}

func (g *Generator) DefaultOutputDir() string {
	return defaultOutputDir
}

func (g *Generator) GenerateThrift() bool {
	return true
}

func (g *Generator) GenerateDependencies(f *parser.Frugal, dir string) error {
	return nil
}

func (g *Generator) GenerateFile(name, outputDir string, fileType generator.FileType) (*os.File, error) {
	if fileType != generator.CombinedFile {
		return nil, fmt.Errorf("frugal: Bad file type for dartlang generator: %s", fileType)
	}
	return g.CreateFile(strings.ToLower(name), outputDir, lang, true)
}

func (g *Generator) GenerateDocStringComment(file *os.File) error {
	comment := fmt.Sprintf(
		"// Autogenerated by Frugal Compiler (%s)\n"+
			"// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING",
		globals.Version)

	_, err := file.WriteString(comment)
	return err
}

func (g *Generator) GenerateServicePackage(file *os.File, f *parser.Frugal, s *parser.Service) error {
	pkg, ok := f.Thrift.Namespaces[lang]
	if ok {
		components := generator.GetPackageComponents(pkg)
		pkg = components[len(components)-1]
	} else {
		pkg = f.Name
	}
	_, err := file.WriteString(fmt.Sprintf("package %s", pkg))
	return err
}

func (g *Generator) GenerateScopePackage(file *os.File, f *parser.Frugal, s *parser.Scope) error {
	pkg, ok := f.Thrift.Namespaces[lang]
	if ok {
		components := generator.GetPackageComponents(pkg)
		pkg = components[len(components)-1]
	} else {
		pkg = f.Name
	}
	_, err := file.WriteString(fmt.Sprintf("package %s", pkg))
	return err
}

func (g *Generator) GenerateServiceImports(file *os.File, s *parser.Service) error {
	imports := "import (\n"
	imports += "\t\"bytes\"\n"
	imports += "\t\"fmt\"\n"
	imports += "\t\"git.apache.org/thrift.git/lib/go/thrift\"\n"
	imports += "\t\"github.com/Workiva/frugal-go\"\n"
	imports += ")\n\n"

	imports += "// (needed to ensure safety because of naive import list construction.)\n"
	imports += "var _ = thrift.ZERO\n"
	imports += "var _ = fmt.Printf\n"
	imports += "var _ = bytes.Equal\n"

	_, err := file.WriteString(imports)
	return err
}

func (g *Generator) GenerateScopeImports(file *os.File, s *parser.Scope) error {
	imports := "import (\n"
	imports += "\t\"fmt\"\n"
	imports += "\t\"log\"\n\n"
	if g.Options["thrift_import"] != "" {
		imports += "\t\"" + g.Options["thrift_import"] + "\"\n"
	} else {
		imports += "\t\"git.apache.org/thrift.git/lib/go/thrift\"\n"
	}
	if g.Options["frugal_import"] != "" {
		imports += "\t\"" + g.Options["frugal_import"] + "\"\n"
	} else {
		imports += "\t\"github.com/Workiva/frugal-go\"\n"
	}

	pkgPrefix := g.Options["package_prefix"]
	for _, include := range f.ReferencedIncludes() {
		namespace, ok := f.NamespaceForInclude(include, lang)
		if !ok {
			namespace = include
		}
		imports += fmt.Sprintf("\t\"%s%s\"\n", pkgPrefix, namespace)
	}

	imports += ")"

	_, err := file.WriteString(imports)
	return err
}

func (g *Generator) GenerateConstants(file *os.File, name string) error {
	if !g.generateConstants {
		return nil
	}
	constants := fmt.Sprintf("const delimiter = \"%s\"", globals.TopicDelimiter)
	_, err := file.WriteString(constants)
	if err != nil {
		return err
	}
	g.generateConstants = false
	return nil
}

func (g *Generator) GeneratePublisher(file *os.File, scope *parser.Scope) error {
	if scope.Comment != nil {
		publishers += g.GenerateInlineComment(scope.Comment, "")
	}
	publishers += fmt.Sprintf("type %sPublisher struct {\n", strings.Title(scope.Name))
	publishers += "\tTransport frugal.Transport\n"
	publishers += "\tProtocol  thrift.TProtocol\n"
	publishers += "\tSeqId     int32\n"
	publishers += "}\n\n"

	publisher += fmt.Sprintf("func New%sPublisher(provider *frugal.Provider) *%sPublisher {\n", strings.Title(scope.Name), strings.Title(scope.Name))
	publisher += "\ttransport, protocol := provider.New()\n"
	publisher += fmt.Sprintf("\treturn &%sPublisher{\n", strings.Title(scope.Name))
	publisher += "\t\tTransport: transport,\n"
	publisher += "\t\tProtocol:  protocol,\n"
	publisher += "\t\tSeqId:     0,\n"
	publisher += "\t}\n"
	publisher += "}\n\n"

	args := ""
	if len(scope.Prefix.Variables) > 0 {
		prefix := ""
		for _, variable := range scope.Prefix.Variables {
			args += prefix + variable
			prefix = ", "
		}
		args += " string, "
	}

	prefix := ""
	for _, op := range scope.Operations {
		publisher += prefix
		prefix = "\n\n"
		if op.Comment != nil {
			publishers += g.GenerateInlineComment(op.Comment, "")
		}
		publishers += fmt.Sprintf("func (l *%sPublisher) Publish%s(%sreq *%s) error {\n",
			strings.Title(scope.Name), op.Name, args, g.qualifiedParamName(op))
		publishers += fmt.Sprintf("\top := \"%s\"\n", op.Name)
		publishers += fmt.Sprintf("\tprefix := %s\n", generatePrefixStringTemplate(scope))
		publishers += "\ttopic := fmt.Sprintf(\"%s" + strings.Title(scope.Name) + "%s%s\", prefix, delimiter, op)\n"
		publishers += "\tl.Transport.PreparePublish(topic)\n"
		publishers += "\toprot := l.Protocol\n"
		publishers += "\tl.SeqId++\n"
		publishers += "\tif err := oprot.WriteMessageBegin(op, thrift.CALL, l.SeqId); err != nil {\n"
		publishers += "\t\treturn err\n"
		publishers += "\t}\n"
		publishers += "\tif err := req.Write(oprot); err != nil {\n"
		publishers += "\t\treturn err\n"
		publishers += "\t}\n"
		publishers += "\tif err := oprot.WriteMessageEnd(); err != nil {\n"
		publishers += "\t\treturn err\n"
		publishers += "\t}\n"
		publishers += "\treturn oprot.Flush()\n"
		publishers += "}"
	}

	_, err := file.WriteString(publisher)
	return err
}

func generatePrefixStringTemplate(scope *parser.Scope) string {
	if len(scope.Prefix.Variables) == 0 {
		if scope.Prefix.String == "" {
			return `""`
		}
		return fmt.Sprintf(`"%s%s"`, scope.Prefix.String, globals.TopicDelimiter)
	}
	template := "fmt.Sprintf(\""
	template += scope.Prefix.Template()
	template += globals.TopicDelimiter + "\", "
	prefix := ""
	for _, variable := range scope.Prefix.Variables {
		template += prefix + variable
		prefix = ", "
	}
	template += ")"
	return template
}

func (g *Generator) GenerateSubscriber(file *os.File, scope *parser.Scope) error {
	if scope.Comment != nil {
		subscribers += g.GenerateInlineComment(scope.Comment, "")
	}
	subscribers += fmt.Sprintf("type %sSubscriber struct {\n", strings.Title(scope.Name))
	subscribers += "\tProvider *frugal.Provider\n"
	subscribers += "}\n\n"

	subscribers += fmt.Sprintf("func New%sSubscriber(provider *frugal.Provider) *%sSubscriber {\n", strings.Title(scope.Name), strings.Title(scope.Name))
	subscribers += fmt.Sprintf("\treturn &%sSubscriber{Provider: provider}\n", strings.Title(scope.Name))
	subscribers += "}\n\n"

	args := ""
	prefix := ""
	if len(scope.Prefix.Variables) > 0 {
		for _, variable := range scope.Prefix.Variables {
			args += prefix + variable
			prefix = ", "
		}
		args += " string, "
	}

	prefix = ""
	for _, op := range scope.Operations {
		subscriber += prefix
		prefix = "\n\n"
		if op.Comment != nil {
			subscribers += g.GenerateInlineComment(op.Comment, "")
		}
		subscribers += fmt.Sprintf("func (l *%sSubscriber) Subscribe%s(%shandler func(*%s)) (*frugal.Subscription, error) {\n",
			strings.Title(scope.Name), op.Name, args, g.qualifiedParamName(op))
		subscribers += fmt.Sprintf("\top := \"%s\"\n", op.Name)
		subscribers += fmt.Sprintf("\tprefix := %s\n", generatePrefixStringTemplate(scope))
		subscribers += "\ttopic := fmt.Sprintf(\"%s" + strings.Title(scope.Name) + "%s%s\", prefix, delimiter, op)\n"
		subscribers += "\ttransport, protocol := l.Provider.New()\n"
		subscribers += "\tif err := transport.Subscribe(topic); err != nil {\n"
		subscribers += "\t\treturn nil, err\n"
		subscribers += "\t}\n\n"
		subscribers += "\tsub := frugal.NewSubscription(topic, transport)\n"
		subscribers += "\tgo func() {\n"
		subscribers += "\t\tfor {\n"
		subscribers += fmt.Sprintf("\t\t\treceived, err := l.recv%s(op, protocol)\n", op.Name)
		subscribers += "\t\t\tif err != nil {\n"
		subscribers += "\t\t\t\tif e, ok := err.(thrift.TTransportException); ok && e.TypeId() == thrift.END_OF_FILE {\n"
		subscribers += "\t\t\t\t\treturn\n"
		subscribers += "\t\t\t\t}\n"
		subscribers += "\t\t\t\tlog.Println(\"frugal: error receiving:\", err)\n"
		subscribers += "\t\t\t\tsub.Signal(err)\n"
		subscribers += "\t\t\t\tsub.Unsubscribe()\n"
		subscribers += "\t\t\t\treturn\n"
		subscribers += "\t\t\t}\n"
		subscribers += "\t\t\thandler(received)\n"
		subscribers += "\t\t}\n"
		subscribers += "\t}()\n\n"
		subscribers += "\treturn sub, nil\n"
		subscribers += "}\n\n"

		subscribers += fmt.Sprintf("func (l *%sSubscriber) recv%s(op string, iprot thrift.TProtocol) (*%s, error) {\n",
			strings.Title(scope.Name), op.Name, g.qualifiedParamName(op))
		subscribers += "\tname, _, _, err := iprot.ReadMessageBegin()\n"
		subscribers += "\tif err != nil {\n"
		subscribers += "\t\treturn nil, err\n"
		subscribers += "\t}\n"
		subscribers += "\tif name != op {\n"
		subscribers += "\t\tiprot.Skip(thrift.STRUCT)\n"
		subscribers += "\t\tiprot.ReadMessageEnd()\n"
		subscribers += "\t\tx9 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, \"Unknown function \"+name)\n"
		subscribers += "\t\treturn nil, x9\n"
		subscribers += "\t}\n"
		subscribers += fmt.Sprintf("\treq := &%s{}\n", g.qualifiedParamName(op))
		subscribers += "\tif err := req.Read(iprot); err != nil {\n"
		subscribers += "\t\treturn nil, err\n"
		subscribers += "\t}\n\n"
		subscribers += "\tiprot.ReadMessageEnd()\n"
		subscribers += "\treturn req, nil\n"
		subscribers += "}"
	}

	_, err := file.WriteString(subscriber)
	return err
}

func (g *Generator) GenerateService(file *os.File, p *parser.Frugal, s *parser.Service) error {
	contents := ""
	contents += g.generateInterface(s)
	contents += g.generateClient(s)
	contents += g.generateServer(s)

	_, err := file.WriteString(contents)
	return err
}

func (g *Generator) generateInterface(service *parser.Service) string {
	contents := fmt.Sprintf("type Frugal%s interface {\n", strings.Title(service.Name))
	for _, method := range service.Methods {
		contents += fmt.Sprintf("\t%s(frugal.Context%s) %s\n",
			strings.Title(method.Name), g.generateInterfaceArgs(method.Arguments),
			g.generateReurnArgs(method))
	}
	contents += "}\n\n"
	return contents
}

func (g *Generator) generateReurnArgs(method *parser.Method) string {
	if method.ReturnType == nil {
		return "(err error)"
	}
	return fmt.Sprintf("(r %s, err error)", g.getGoTypeFromThriftType(method.ReturnType))
}

func (g *Generator) generateClient(service *parser.Service) string {
	servTitle := strings.Title(service.Name)

	contents := fmt.Sprintf("type Frugal%sClient struct {\n", servTitle)
	contents += "\tTransport       thrift.TTransport\n"
	contents += "\tProtocolFactory frugal.FProtocolFactory\n"
	contents += "\tInputProtocol   frugal.FProtocol\n"
	contents += "\tOutputProtocol  frugal.FProtocol\n"
	contents += "\tSeqId           int32\n"
	contents += "}\n\n"

	contents += fmt.Sprintf(
		"func NewFrugal%sClientFactory(t thrift.TTransport, f frugal.FProtocolFactory) *Frugal%sClient {\n",
		servTitle, servTitle)
	contents += fmt.Sprintf("\treturn &Frugal%sClient{\n", servTitle)
	contents += "\t\tTransport:       t,\n"
	contents += "\t\tProtocolFactory: f,\n"
	contents += "\t\tInputProtocol:   f.GetProtocol(t),\n"
	contents += "\t\tOutputProtocol:  f.GetProtocol(t),\n"
	contents += "\t\tSeqId:           0,\n"
	contents += "\t}\n"
	contents += "}\n\n"

	contents += fmt.Sprintf(
		"func NewFrugal%sClientProtocol(t thrift.TTransport, iprot, oprot frugal.FProtocol) *Frugal%sClient {\n",
		service.Name, service.Name)
	contents += fmt.Sprintf("\treturn &Frugal%sClient{\n", servTitle)
	contents += "\t\tTransport:       t,\n"
	contents += "\t\tProtocolFactory: nil,\n"
	contents += "\t\tInputProtocol:   iprot,\n"
	contents += "\t\tOutputProtocol:  oprot,\n"
	contents += "\t\tSeqId:           0,\n"
	contents += "\t}\n"
	contents += "}\n\n"

	for _, method := range service.Methods {
		contents += g.generateClientMethod(service, method)
	}
	return contents
}

func (g *Generator) generateClientMethod(service *parser.Service, method *parser.Method) string {
	servTitle := strings.Title(service.Name)
	nameTitle := strings.Title(method.Name)
	nameLower := strings.ToLower(method.Name)

	contents := fmt.Sprintf("func (f *Frugal%sClient) %s(ctx frugal.Context%s) %s {\n",
		servTitle, nameTitle, g.generateInputArgs(method.Arguments),
		g.generateReurnArgs(method))
	contents += fmt.Sprintf("\tif err = f.send%s(ctx, %s); err != nil {\n",
		nameTitle, g.generateClientOutputArgs(method.Arguments))
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += fmt.Sprintf("\treturn f.recv%s(ctx)\n", nameTitle)
	contents += "}\n\n"

	contents += fmt.Sprintf("func (f *Frugal%sClient) send%s(ctx frugal.Context%s) (err error) {\n",
		servTitle, nameTitle, g.generateInputArgs(method.Arguments))
	contents += "\toprot := f.OutputProtocol\n"
	contents += "\tif oprot == nil {\n"
	contents += "\t\toprot = f.ProtocolFactory.GetProtocol(f.Transport)\n"
	contents += "\t\tf.OutputProtocol = oprot\n"
	contents += "\t}\n"
	contents += fmt.Sprintf("\tif err = f.OutputProtocol.WriteRequestHeader(ctx); err != nil {\n")
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += "\tf.SeqId++\n"
	contents += fmt.Sprintf(
		"\tif err = oprot.WriteMessageBegin(\"%s\", thrift.CALL, f.SeqId); err != nil {\n",
		nameLower)
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += fmt.Sprintf("\targs := %s%sArgs{\n", servTitle, nameTitle)
	contents += g.generateStructArgs(method.Arguments)
	contents += "\t}\n"
	contents += "\tif err = args.Write(oprot); err != nil {\n"
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += "\tif err = oprot.WriteMessageEnd(); err != nil {\n"
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += "\treturn oprot.Flush()\n"
	contents += "}\n\n"

	contents += fmt.Sprintf("func (f *Frugal%sClient) recv%s(ctx frugal.Context) %s {\n",
		servTitle, nameTitle, g.generateReurnArgs(method))
	contents += "\tiprot := f.InputProtocol\n"
	contents += "\tif iprot == nil {\n"
	contents += "\t\tiprot = f.ProtocolFactory.GetProtocol(f.Transport)\n"
	contents += "\t\tf.InputProtocol = iprot\n"
	contents += "\t}\n"
	contents += "\tif err = iprot.ReadResponseHeader(ctx); err != nil {\n"
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += "\tmethod, mTypeId, seqId, err := iprot.ReadMessageBegin()\n"
	contents += "\tif err != nil {\n"
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += fmt.Sprintf("\tif method != \"%s\" {\n", nameLower)
	contents += fmt.Sprintf(
		"\terr = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, \"%s failed: wrong method name\")\n",
		nameLower)
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += "\tif f.SeqId != seqId {\n"
	contents += fmt.Sprintf(
		"\terr = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, \"%s failed: out of sequence response\")\n",
		nameLower)
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += "\tif mTypeId == thrift.EXCEPTION {\n"
	contents += "\t\terror0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, \"Unknown Exception\")\n"
	contents += "\t\tvar error1 error\n"
	contents += "\t\terror1, err = error0.Read(iprot)\n"
	contents += "\t\tif err != nil {\n"
	contents += "\t\t\treturn\n"
	contents += "\t\t}\n"
	contents += "\t\tif err = iprot.ReadMessageEnd(); err != nil {\n"
	contents += "\t\t\treturn\n"
	contents += "\t\t}\n"
	contents += "\t\terr = error1\n"
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += "\tif mTypeId != thrift.REPLY {\n"
	contents += fmt.Sprintf(
		"\terr = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, \"%s failed: invalid message type\")\n",
		nameLower)
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += fmt.Sprintf("\tresult := %s%sResult{}\n", servTitle, nameTitle)
	contents += "\tif err = result.Read(iprot); err != nil {\n"
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += "\tif err = iprot.ReadMessageEnd(); err != nil {\n"
	contents += "\t\treturn\n"
	contents += "\t}\n"
	for _, err := range method.Exceptions {
		errTitle := strings.Title(err.Name)
		contents += fmt.Sprintf("\tif result.%s != nil {\n", errTitle)
		contents += fmt.Sprintf("\t\terr = result.%s\n", errTitle)
		contents += "\t\treturn\n"
		contents += "\t}\n"
	}
	if method.ReturnType != nil {
		contents += "\tr = result.GetSuccess()\n"
	}
	contents += "\treturn\n"
	contents += "}\n\n"

	return contents
}

func (g *Generator) generateServer(service *parser.Service) string {
	contents := ""
	contents += g.generateProcessor(service)
	for _, method := range service.Methods {
		contents += g.generateMethodProcessor(service, method)
	}
	return contents
}

func (g *Generator) generateProcessor(service *parser.Service) string {
	servTitle := strings.Title(service.Name)
	servLower := strings.ToLower(service.Name)
	contents := ""
	contents += fmt.Sprintf("type Frugal%sProcessor struct {\n", servTitle)
	contents += "\tprocessorMap map[string]frugal.FProcessorFunction\n"
	contents += fmt.Sprintf("\thandler      Frugal%s\n", servTitle)
	contents += "}\n\n"

	contents += fmt.Sprintf("func (p *Frugal%sProcessor) GetProcessorFunction(key string) "+
		"(processor frugal.FProcessorFunction, ok bool) {\n", servTitle)
	contents += "\tprocessor, ok = p.processorMap[key]\n"
	contents += "\treturn\n"
	contents += "}\n\n"

	contents += fmt.Sprintf("func NewFrugal%sProcessor(handler Frugal%s) *Frugal%sProcessor {\n",
		servTitle, servTitle, servTitle)
	contents += fmt.Sprintf("\tp := &Frugal%sProcessor{\n", servTitle)
	contents += "\t\thandler:      handler,\n"
	contents += "\t\tprocessorMap: make(map[string]frugal.FProcessorFunction),\n"
	contents += "\t}\n"
	for _, method := range service.Methods {
		contents += fmt.Sprintf("\tp.processorMap[\"%s\"] = &%sFrugalProcessor%s{handler: handler}\n",
			strings.ToLower(method.Name), servLower, strings.Title(method.Name))
	}
	contents += "\treturn p\n"
	contents += "}\n\n"

	contents += fmt.Sprintf("func (p *Frugal%sProcessor) Process(iprot, oprot frugal.FProtocol) "+
		"(success bool, err thrift.TException) {\n", servTitle)
	contents += "\tctx, err := iprot.ReadRequestHeader()\n"
	contents += "\tif err != nil {\n"
	contents += "\t\treturn false, err\n"
	contents += "\t}\n"
	contents += "\tname, _, seqId, err := iprot.ReadMessageBegin()\n"
	contents += "\tif err != nil {\n"
	contents += "\t\treturn false, err\n"
	contents += "\t}\n"
	contents += "\tif processor, ok := p.GetProcessorFunction(name); ok {\n"
	contents += "\t\treturn processor.Process(ctx, seqId, iprot, oprot)\n"
	contents += "\t}\n"
	contents += "\tiprot.Skip(thrift.STRUCT)\n"
	contents += "\tiprot.ReadMessageEnd()\n"
	contents += "\tx3 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, \"Unknown function \"+name)\n"
	contents += "\toprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)\n"
	contents += "\tx3.Write(oprot)\n"
	contents += "\toprot.WriteMessageEnd()\n"
	contents += "\toprot.Flush()\n"
	contents += "\treturn false, x3\n"
	contents += "}\n\n"

	return contents
}

func (g *Generator) generateMethodProcessor(service *parser.Service, method *parser.Method) string {
	servTitle := strings.Title(service.Name)
	servLower := strings.ToLower(service.Name)
	nameTitle := strings.Title(method.Name)
	nameLower := strings.ToLower(method.Name)

	contents := fmt.Sprintf("type %sFrugalProcessor%s struct {\n", servLower, nameTitle)
	contents += fmt.Sprintf("\thandler Frugal%s\n", servTitle)
	contents += "}\n\n"

	contents += fmt.Sprintf("func (p *%sFrugalProcessor%s) Process(ctx frugal.Context, "+
		"seqId int32, iprot, oprot frugal.FProtocol) (success bool, err thrift.TException) {\n",
		servLower, nameTitle)
	contents += fmt.Sprintf("\targs := %s%sArgs{}\n", servTitle, nameTitle)
	contents += "\tif err = args.Read(iprot); err != nil {\n"
	contents += "\t\tiprot.ReadMessageEnd()\n"
	contents += "\t\tx := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())\n"
	contents += fmt.Sprintf("\t\toprot.WriteMessageBegin(\"%s\", thrift.EXCEPTION, seqId)\n",
		nameLower)
	contents += "\t\tx.Write(oprot)\n"
	contents += "\t\toprot.WriteMessageEnd()\n"
	contents += "\t\toprot.Flush()\n"
	contents += "\t\treturn false, err\n"
	contents += "\t}\n\n"

	contents += "\tiprot.ReadMessageEnd()\n"
	contents += fmt.Sprintf("\tresult := %s%sResult{}\n", servTitle, nameTitle)
	contents += "\tvar err2 error\n"
	if method.ReturnType != nil {
		contents += fmt.Sprintf("\tvar retval %s\n", g.getGoTypeFromThriftType(method.ReturnType))
		contents += fmt.Sprintf("\tif retval, err2 = p.handler.%s(ctx, %s); err2 != nil {\n",
			nameTitle, g.generateServerOutputArgs(method.Arguments))
	} else {
		contents += fmt.Sprintf("\tif err2 = p.handler.%s(ctx, %s); err2 != nil {\n",
			nameTitle, g.generateServerOutputArgs(method.Arguments))
	}
	if len(method.Exceptions) > 0 {
		contents += "\t\tswitch v := err2.(type) {\n"
		for _, err := range method.Exceptions {
			contents += fmt.Sprintf("\t\tcase %s:\n", g.getGoTypeFromThriftType(err.Type))
			contents += fmt.Sprintf("\t\t\tresult.%s = v\n", strings.Title(err.Name))
		}
		contents += "\t\tdefault:\n"
		contents += g.generateMethodException("\t\t\t", nameLower)
		contents += "\t\t}\n"
	} else {
		contents += g.generateMethodException("\t\t", nameLower)
	}
	if method.ReturnType != nil {
		contents += "\t} else {\n"
		contents += "\t\tresult.Success = &retval\n"
	}
	contents += "\t}\n"

	contents += "\tif err2 = oprot.WriteResponseHeader(ctx); err2 != nil {\n"
	contents += "\t\terr = err2\n"
	contents += "\t}\n"
	contents += fmt.Sprintf("\tif err2 = oprot.WriteMessageBegin(\"%s\", "+
		"thrift.REPLY, seqId); err2 != nil {\n", nameLower)
	contents += "\t\terr = err2\n"
	contents += "\t}\n"
	contents += "\tif err2 = result.Write(oprot); err == nil && err2 != nil {\n"
	contents += "\t\terr = err2\n"
	contents += "\t}\n"
	contents += "\tif err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {\n"
	contents += "\t\terr = err2\n"
	contents += "\t}\n"
	contents += "\tif err2 = oprot.Flush(); err == nil && err2 != nil {\n"
	contents += "\t\terr = err2\n"
	contents += "\t}\n"
	contents += "\tif err != nil {\n"
	contents += "\t\treturn\n"
	contents += "\t}\n"
	contents += "\treturn true, err\n"
	contents += "}\n\n"

	return contents
}

func (g *Generator) generateMethodException(prefix, methodName string) string {
	contents := fmt.Sprintf(prefix+"x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "+
		"\"Internal error processing %s: \"+err2.Error())\n", methodName)
	contents += fmt.Sprintf(prefix+"oprot.WriteMessageBegin(\"%s\", thrift.EXCEPTION, seqId)\n",
		methodName)
	contents += prefix + "x.Write(oprot)\n"
	contents += prefix + "oprot.WriteMessageEnd()\n"
	contents += prefix + "oprot.Flush()\n"
	contents += prefix + "return true, err2\n"
	return contents
}

func (g *Generator) generateInterfaceArgs(args []*parser.Field) string {
	argStr := ""
	for _, arg := range args {
		argStr += ", " + g.getGoTypeFromThriftType(arg.Type)
	}
	return argStr
}

func (g *Generator) generateClientOutputArgs(args []*parser.Field) string {
	argStr := ""
	for i, arg := range args {
		argStr += arg.Name
		if i < len(args)-1 {
			argStr += ", "
		}
	}
	return argStr
}

func (g *Generator) generateInputArgs(args []*parser.Field) string {
	argStr := ""
	for _, arg := range args {
		argStr += ", " + arg.Name + " " + g.getGoTypeFromThriftType(arg.Type)
	}
	return argStr
}

func (g *Generator) generateStructArgs(args []*parser.Field) string {
	argStr := ""
	for _, arg := range args {
		argStr += "\t\t" + strings.Title(arg.Name) + ": " + arg.Name + ",\n"
	}
	return argStr
}

func (g *Generator) generateServerOutputArgs(args []*parser.Field) string {
	argStr := ""
	for i, arg := range args {
		argStr += fmt.Sprintf("args.%s", strings.Title(arg.Name))
		if i < len(args)-1 {
			argStr += ", "
		}
	}
	return argStr
}

func (g *Generator) getGoTypeFromThriftType(t *parser.Type) string {
	switch t.Name {
	case "bool":
		return "bool"
	case "byte":
		return "byte"
	case "i16":
		return "int16"
	case "i32":
		return "int32"
	case "i64":
		return "int64"
	case "double":
		return "float64"
	case "string":
		return "string"
	case "binary":
		return "[]byte"
	case "list":
		return fmt.Sprintf("[]%s", g.getGoTypeFromThriftType(t.ValueType))
	case "set":
		return fmt.Sprintf("map[%s]bool", g.getGoTypeFromThriftType(t.ValueType))
	case "map":
		return fmt.Sprintf("map[%s]%s", g.getGoTypeFromThriftType(t.KeyType),
			g.getGoTypeFromThriftType(t.ValueType))
	default:
		// This is a custom type, return a pointer to it
		return "*" + t.Name
	}
}

func (g *Generator) qualifiedParamName(op *parser.Operation) string {
	param := op.ParamName()
	include := op.IncludeName()
	if include != "" {
		namespace, ok := g.Frugal.NamespaceForInclude(include, lang)
		if !ok {
			namespace = include
		}
		param = fmt.Sprintf("%s.%s", namespace, param)
	}
	return param
}

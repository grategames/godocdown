package main

import (
	"fmt"
	"go/doc"
	"io"
)

func renderConstantSectionTo(writer io.Writer, list []*doc.Value) {
	for _, entry := range list {
		fmt.Fprintf(writer, "%s\n%s\n", indentCode(sourceOfNode(entry.Decl)), formatIndent(filterText(entry.Doc)))
	}
}

func renderVariableSectionTo(writer io.Writer, list []*doc.Value) {
	for _, entry := range list {
		fmt.Fprintf(writer, "%s\n%s\n", indentCode(sourceOfNode(entry.Decl)), formatIndent(filterText(entry.Doc)))
	}
}

func renderFunctionSectionTo(writer io.Writer, list []*doc.Func, inTypeSection bool) {

	header := RenderStyle.FunctionHeader
	if inTypeSection {
		header = RenderStyle.TypeFunctionHeader
	}

	for _, entry := range list {
		receiver := " "
		if entry.Recv != "" {
			receiver = fmt.Sprintf("(%s) ", entry.Recv)
		}
		fmt.Fprintf(writer, "%s func %s%s\n\n%s\n%s\n", header, receiver, entry.Name, indentCode(sourceOfNode(entry.Decl)), formatIndent(filterText(entry.Doc)))
	}
}

func renderLinksSectionTo(writer io.Writer, doc *_document) {
	counter := 1
	for _, entry := range doc.pkg.Funcs {
		counter++
		fmt.Fprintf(writer, "[func  %s](#toc_%s)  \n", entry.Name, fmt.Sprint(counter))
	}
	for _, entry := range doc.pkg.Types {
		counter++
		fmt.Fprintf(writer, "[type  %s](#toc_%s)  \n", entry.Name, fmt.Sprint(counter))
		
		for _, entry := range entry.Funcs {
			counter++
			fmt.Fprintf(writer, "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [func  %s](#toc_%s)  \n", entry.Name, fmt.Sprint(counter))
		}
		
		for _, entry := range entry.Methods {
			counter++
			receiver := " "
			if entry.Recv != "" {
				receiver = fmt.Sprintf("(%s) ", entry.Recv)
			}
			fmt.Fprintf(writer, "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; [func %s%s](#toc_%s)  \n", receiver, entry.Name, fmt.Sprint(counter))
		}
	}
	fmt.Fprint(writer, "\n")
}

func renderTypeSectionTo(writer io.Writer, list []*doc.Type) {

	header := RenderStyle.TypeHeader

	for _, entry := range list {
		fmt.Fprintf(writer, "%s type %s\n\n%s\n\n%s\n", header, entry.Name, indentCode(sourceOfNode(entry.Decl)), formatIndent(filterText(entry.Doc)))
		renderConstantSectionTo(writer, entry.Consts)
		renderVariableSectionTo(writer, entry.Vars)
		renderFunctionSectionTo(writer, entry.Funcs, true)
		renderFunctionSectionTo(writer, entry.Methods, true)
	}
}

func renderHeaderTo(writer io.Writer, document *_document) {
	fmt.Fprintf(writer, "# %s\n--\n", document.Name)

	if !document.IsCommand {
		// Import
		if RenderStyle.IncludeImport {
			if document.ImportPath != "" {
				fmt.Fprintf(writer, spacer(4)+"import \"%s\"\n\n", document.ImportPath)
			}
		}
	}
}

func renderSynopsisTo(writer io.Writer, document *_document) {
	fmt.Fprintf(writer, "%s\n", headifySynopsis(formatIndent(filterText(document.pkg.Doc))))
}

func renderUsageTo(writer io.Writer, document *_document) {
	// Links
	renderLinksSectionTo(writer, document)

	// Usage
	fmt.Fprintf(writer, "%s\n", RenderStyle.UsageHeader)

	// Constant Section
	renderConstantSectionTo(writer, document.pkg.Consts)

	// Variable Section
	renderVariableSectionTo(writer, document.pkg.Vars)

	// Function Section
	renderFunctionSectionTo(writer, document.pkg.Funcs, false)

	// Type Section
	renderTypeSectionTo(writer, document.pkg.Types)
}

func renderSignatureTo(writer io.Writer) {
	if RenderStyle.IncludeSignature {
		fmt.Fprintf(writer, "\n\n--\n**godocdown** http://github.com/robertkrimen/godocdown\n")
	}
}

package main

import (
	"fmt"
	"net/http"
	"time"

	"strings"

	"github.com/jung-kurt/gofpdf"
)

const assetDir = "assets/"

func main() {
	fs := http.FileServer(http.Dir(assetDir))
	http.Handle("/", fs)
	http.HandleFunc("/api/render-pdf", renderPdf)

	http.ListenAndServe(":8080", nil)
}

func renderPdf(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse form data: %v", err), http.StatusBadRequest)
		return
	}
	htmlName := assetDir[:len(assetDir)-1] + r.FormValue("_referrer")
	content, err := getPrintableContent(htmlName)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get printable content for %v: %v", htmlName, err), http.StatusBadRequest)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 16)
	printParagraphs(content, pdf)
	printDate(pdf)
	printForm(content, r, pdf)
	printSignature(r, pdf)
	err = pdf.Output(w)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not create pdf: %v", err), http.StatusInternalServerError)
		return
	}
}

func printDate(pdf *gofpdf.Fpdf) {
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Local().Format("Mon Jan 2 2006")))
	pdf.Ln(-1)
}

func printParagraphs(content *printable, pdf *gofpdf.Fpdf) {
	_, lineHt := pdf.GetFontSize()
	html := pdf.HTMLBasicNew()
	for _, par := range content.paragraphs {
		html.Write(lineHt, par)
		// Line break
		pdf.Ln(lineHt)
		pdf.Ln(lineHt)
	}
}

func printForm(content *printable, r *http.Request, pdf *gofpdf.Fpdf) {
	for k, v := range r.Form {
		if !strings.HasPrefix(k, "_") {
			for _, str := range v {
				pdf.Cell(40, 10, fmt.Sprintf("%s: %s", k, str))
				pdf.Ln(-1)
			}
		}
	}
}

func printSignature(r *http.Request, pdf *gofpdf.Fpdf) {
	svg := r.FormValue("_sigval")
	sig, err := gofpdf.SVGBasicParse([]byte(svg))
	if err == nil {
		scale := 100 / sig.Wd
		scaleY := 30 / sig.Ht
		if scale > scaleY {
			scale = scaleY
		}
		pdf.SetLineCapStyle("round")
		pdf.SetLineWidth(0.25)
		pdf.SetY(pdf.GetY() + 10)
		pdf.SVGBasicWrite(&sig, scale)
	} else {
		pdf.SetError(err)
	}
}

#import "@preview/i-figured:0.2.4"


#let thesis_template(content) = {
  set page(
    paper: "a4",
    margin: (
      left: 3cm,
      right: 2cm,
      y: 2cm,
    ),
    columns: 1,
    number-align: bottom + center,
  )
  
  set heading(
    numbering: "1.",
  )
  
  show heading: it => {
    set text(
      size: 14pt,
      weight: "bold",
    )
    
    it
  }
  
  set par(
    first-line-indent: 1.25cm,
    leading: 1.5em,
    justify: true,
  )
  
  set text(
    font: "Times New Roman",
    size: 12pt,
    lang: "ru",
  )
  
  set figure(
    supplement: auto,
    placement: auto,
  )
  
  show figure.caption: emph

  // show heading: i-figured.reset-counters
  show figure: it => i-figured.show-figure(
    it,
    level: 0,
    numbering: "1",
  )

  set list(
    indent: 0.75cm,
  )

  set enum(
    indent: 0.75cm,
  )

  // show " .": it => highlight(fill: red, it)
  // show " ,": it => highlight(fill: red, it)

  content
}


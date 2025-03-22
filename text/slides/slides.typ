#import "/utils.typ": *
#import "@preview/polylux:0.4.0": *
#import "./metropolis-theme.typ" as metropolis

#set document(
  author: "Беляев Андрей Алексеевич",
  title: [
    #todo[]
  ]
)

#show: metropolis.setup.with()

#set text(
  lang: "ru",
)

#let center-img(
  path,
  width: 100%,
  height: 100%,
  fit: "contain",
) = align(
  center + horizon,
  image(
    path,
    width: width,
    height: height,
    fit: fit,
  )
)

// #hide_review_notes.update(false)


#slide[
  #set page(header: none, footer: none, margin: 3em)
  #set align(horizon)
  
  #text(size: 1.3em)[
    Резервное копирование облачного хранилища Microsoft Azure Blob Storage
  ]

  #metropolis.divider
  
  #set text(size: .8em, weight: "light")
  Беляев Андрей Алексеевич, М05-411

  #text(style: "italic")[
    #todo[Дата]
  ]
]

#slide[
  = Введение
  
  #todo[]
]

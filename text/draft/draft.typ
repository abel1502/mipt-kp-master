#import "/utils.typ": *

#let paper(
  name: "",
  file: "",
  text,
) = [
  == #name
  _(#file)_

  #text
]

= Обзор литературы (набросок)

// deduplication

#paper(
  name: "Building a High-performance Deduplication System",
  file: "GuoEfstathopoulos.pdf",
)[
  #todo[]
]

#paper(
  name: "Demystifying Data Deduplication",
  file: "mandagere2008.pdf",
)[
  #todo[]
]

#paper(
  name: "A Study of Practical Deduplication",
  file: "meyer2012.pdf",
)[
  #todo[]
]

#paper(
  name: "Data Deduplication techniques",
  file: "qinluhe2010.pdf",
)[
  #todo[]
]

#paper(
  name: "A Comprehensive Study of thePast, Present, and Future of
Data Deduplication",
  file: "xia2016.pdf",
)[
  #todo[]
]

// encryption

#paper(
  name: "A Study on Deduplication Techniques over Encrypted Data",
  file: "akhila2016.pdf",
)[
  #todo[]
]

#paper(
  name: "DupLESS: Server-Aided Encryption for Deduplicated Storage",
  file: "sec13-paper_bellare.pdf",
)[
  #todo[]
]

#paper(
  name: "Encrypted Data Management with Deduplication in Cloud Computing",
  file: "yan2016-2.pdf",
)[
  #todo[]
]

#paper(
  name: "Deduplication on Encrypted Big Data in Cloud",
  file: "yan2016.pdf",
)[
  #todo[]
]


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
  Вводятся три механизма, призванных улучшить эффективность дедупликации:
  
  + Прогрессивное выборочное индексирование (_Progressive sampled indexing_). Индекс --- сопоставление идентификатора фрагмента (хеша) и расположения его полного содержимого. Вместо единого индекса, содержащего информацию о всех блоках, предлагается хранить информацию, актуальную для конкретного файла, вместе с остальными метаданными файла (в списке его блоков). В таком случае вместо полного глобального индекса можно использовать выборочный, содержащий информацию только о некоторых "горячих" фрагментах. Предлагаются некоторые методы спекулятивного принятия решений о взятии фрагментов в индекс. Прогрессивность относится к динамическому выбору размера "горячего" индекса (_sampling rate_) в обратной пропорции к объёму хранимых данных.

  + Сгруппированная пометка-и-уборка (_Grouped mark-and-sweep_). #todo[Подробнее]

  + Многопоточная модель взаимодействия клиент-сервер на основе событий (_Event-driven multi-threaded client-server interaction model_). #todo[Подробнее]

  Благодаря предложенным механизмам разработанный прототип продемонстрировал отличную масштабируемость, высокую пропускную способность и низкую деградацию эффективности дедупликации.
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


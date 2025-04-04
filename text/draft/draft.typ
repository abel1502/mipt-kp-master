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
  
  + Прогрессивное выборочное индексирование (_Progressive sampled indexing_).
    Индекс --- сопоставление идентификатора фрагмента (хеша) и расположения его полного содержимого. Вместо единого индекса, содержащего информацию о всех блоках, предлагается хранить информацию, актуальную для конкретного файла, вместе с остальными метаданными файла (в списке его блоков). В таком случае вместо полного глобального индекса можно использовать выборочный, содержащий информацию только о некоторых "горячих" фрагментах. Предлагаются некоторые методы спекулятивного принятия решений о взятии фрагментов в индекс. Прогрессивность относится к динамическому выбору размера "горячего" индекса (_sampling rate_) в обратной пропорции к объёму хранимых данных.

  + Сгруппированная пометка-и-уборка (_Grouped mark-and-sweep_).
    #todo[Подробнее]

  + Многопоточная модель взаимодействия клиент-сервер на основе событий (_Event-driven multi-threaded client-server interaction model_).
    #todo[Подробнее]

  Благодаря предложенным механизмам разработанный прототип продемонстрировал отличную масштабируемость, высокую пропускную способность и низкую деградацию эффективности дедупликации.
]

#paper(
  name: "Demystifying Data Deduplication",
  file: "mandagere2008.pdf",
)[
  Обозреваются таксономия понятий в дедупликации и численные параметры, сопутствующие различным вариантам резервного копирования, в реальных условиях.

  Выделяются три свойства решений в дедупликации:

  + Размещение функционала дедупликации
    + Клиент-сервер резервного копирования.
    + Дедуплицирующая аппаратура.
    + Массив хранения данных.

  + Время проведения дедупликации
    + Синхронная дедупликация (_In-Band_).
    + Асинхронная дедупликация (_Out-of-band_).

  + Алгоритм поиска повторов
    + Дельта-кодирование.
    + Хеширование целых файлов.
    + Хеширование блоков фиксированного размера.
    + Хеширование блоков переменного размера.

  #todo[Подробнее; ещё?]
]

#paper(
  name: "A Study of Practical Deduplication",
  file: "meyer2012.pdf",
)[
  Для персональных компьютеров была проведена оценка сравнительной эффективности различных методов дедупликации. Было получено, что дедупликация на уровне целых файлов достигает 75% экономии места, обеспечиваемой наиболее агрессивной дедупликацией на уровне отдельных блоков, при использовании в активных файловых системах, и 87% при использовании в резервных копиях. Также было получено, что распределение размеров файлов смещается в направлении крупных неструктурированных файлов, а фрагментация данных на диске
  на практике несущественна.
]

#paper(
  name: "Data Deduplication techniques",
  file: "qinluhe2010.pdf",
)[
  _Работа, при более детальном изучении, не заслуживает рассмотрения в НИР._
]

#paper(
  name: "A Comprehensive Study of thePast, Present, and Future of
Data Deduplication",
  file: "xia2016.pdf",
)[
  Рассмотрена таксономия понятий в дедупликации и достижения, описанные в других работах, а также перспективные открытые проблемы.

  Приведена статистика степени сжатия данных за счёт дедупликации в различных практических контекстах --- на персональных компьютерах и на серверах, при блочной или полнофайловой дедупликации, и других параметрах --- полученная на основе 5 других работ. Рассмотрены механизмы разбиения потока данных на блоки, основанные на алгоритме Рабина, с различными доработками, предложенными в 13 других работах. Также рассмотрены способы вычислительной оптимизации задач дедупликации, предложенные в 7 других работах. Кроме того, рассмотрены способы индексирования известных блоков и поиска сходств / дельта-кодирования. Проведён краткий обзор подходов ко внедрению криптографии в системы с дедупликацией. Рассмотрены сценарии, где дедупликация может приносить практическую пользу.

  _Статья заслуживает большого внимания в НИР._
]

// encryption

#paper(
  name: "A Study on Deduplication Techniques over Encrypted Data",
  file: "akhila2016.pdf",
)[
  Проведён обзор методов использования шифрования совместно с дедупликацией, предложенных в других работах. Можно выделить принципиальные методы:
  - Шифрование, завязанное на сообщение (_message-locked encryption_) / конвергентное шифрование (_convergent encryption_). Ключ шифрования зависит только от содержимого шифруемого сообщения. Позволяет сравнивать зашифрованные сообщения точно так же, как и исходные.
  - Доказательство владения (_proof of ownership_). Интерактивный алгоритм, подтверждающий владение полными данными, соответствующими некоторому хешу. Необходимо для противодействия атакам внедрения фиктивных данных, якобы дублирующих некоторый секретный блок.
  - Генерация ключей с участием сервера. Многие методы полагаются на доверенную третью сторону, участвующую в генерации ключей. Это позволяет реализовать конвергентное шифрование, не допуская при этом атак по раскрытию открытого текста из узкого спектра.
  - Хранение "популярных" данных без шифрования. Предполагается, что данные с высокой степенью дупликации не носят секретный характер и могут храниться без шифрования.
]

#paper(
  name: "DupLESS: Server-Aided Encryption for Deduplicated Storage",
  file: "sec13-paper_bellare.pdf",
)[
  Предлагается способ усиления шифрования, завязанного на сообщения, призванный устранить возможность атак, раскрывающих открытый текст из узкого спектра. При этом используется сервер ключей. Предоставлено описание криптографического протокола и его реализация, работающая с существующими крупными облачными хранилищами без какой-либо особой поддержки с их стороны. Разработанный протокол не приводит к значительным потерям эффективности дедупликации. Безопасность протокола достигается за счёт использования дополнительного секретного параметра при генерации ключей, а также использования протокола забывчивой передачи, чтобы избежать разглашения этого параметра.
]

#paper(
  name: "Encrypted Data Management with Deduplication in Cloud Computing",
  file: "yan2016-2.pdf",
)[
  Предлагается способ использования шифрования с дедупликацией за счёт шифрования на основе аттрибутов (_attribute-based encryption_).

  #todo[]
]

#paper(
  name: "Deduplication on Encrypted Big Data in Cloud",
  file: "yan2016.pdf",
)[
  Предлагается способ использования шифрования с дедупликацией за счёт проверки владения и опосредованного пере-шифрования (_ownership challenge and proxy re-encryption_).

  #todo[]
]


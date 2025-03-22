#let nullpar = par(h(0pt))


#let with(ctx_var, value, cont) = context {
  let old_value = ctx_var.get() 
  ctx_var.update(value)
  cont
  ctx_var.update(old_value)
}


#let hide_review_notes = state("hide_review_notes", false)

#let todo(cont) = {
  if cont != [] {
    cont = [: ] + cont
  }
  context if hide_review_notes.get() {
    []
  } else {
    highlight(fill: yellow)[[TODO#cont]]
  }
}

#let needs_review(cont) = {
  context if hide_review_notes.get() {
    cont
  } else {
    highlight(fill: yellow)[#cont \[?\]]
  }
}


#let edits_display = state("edits_display", "both")

#let cross(cont) = context (:
  old: cont,
  new: [],
  both: {
    set text(fill: red)
    strike(cont)
  },
).at(edits_display.get())

#let add(cont) = context (:
  old: [],
  new: cont,
  both: {
    set text(fill: green)
    cont
  },
).at(edits_display.get())

#let replace(cont, repl) = context (:
  old: cont,
  new: repl,
  both: {
    cross(cont)
    [ ]
    add(repl)
  },
).at(edits_display.get())


#let comment(cont) = {
  set text(fill: green.darken(50%))
  [\/\/ ]
  cont
  [\ ]
}

#let comments(..cont) = {
  set par(first-line-indent: 0em)

  nullpar

  cont.pos().map(comment).join()
}

#let hide(cont) = []

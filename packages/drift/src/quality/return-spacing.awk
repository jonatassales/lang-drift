# Prints, once, each file with a return that follows another statement of its
# block without a blank line before it (.claude/rules/00-general.rule.md,
# Statement spacing). A return that opens its block, follows a case label or
# is the body of an if, else, for or while without braces needs none; a comment
# right above the return counts as part of it.
FNR == 1 {
  last = ""
  gap = 0
  guard = 0
  header = 0
}

{
  line = $0
  sub(/^[ \t]+/, "", line)
}

line == "" {
  gap = 1
  next
}

line ~ /^(\/\/|\/\*|\*)/ {
  next
}

line ~ /^return([ \t;(]|$)/ && !gap && last != "" && last !~ /[{(\[:]$/ && !guard && !(FILENAME in reported) {
  print FILENAME
  reported[FILENAME] = 1
}

# A guard's condition may span lines: it starts at if, for or while and ends
# where its parentheses balance, without a brace or a semicolon.
line ~ /^(\}[ \t]*)?(if|else|for|while)([ \t(]|$)/ {
  header = 1
  depth = 0
}

{
  opened = gsub(/\(/, "(", line)
  closed = gsub(/\)/, ")", line)
  depth += opened - closed
  guard = header && depth <= 0 && line !~ /[;{}]$/
  if (depth <= 0) header = 0
  last = line
  gap = 0
}

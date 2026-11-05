# Prints, once, each file with an await statement that has no blank line before
# or after it (.claude/rules/04-typescript.rule.md, Statement spacing). The
# first and the last statement of a block need none, an await inside a nested
# function does not count, a comment above a statement belongs to it, and a
# statement that spans lines is judged from where it starts to where its
# brackets close.
FNR == 1 {
  previous = ""
  pending = 0
  statement = 0
}

{
  line = $0
  sub(/^[ \t]+/, "", line)
}

# The line after an await statement that ended: blank, a closing bracket or a
# continuation of the same expression.
pending {
  if (line != "" && line !~ /^([})\]]|\.|\?|:|&&|\|\|)/ && !(FILENAME in reported)) {
    print FILENAME
    reported[FILENAME] = 1
  }

  pending = 0
}

line ~ /^(\/\/|\/\*|\*)/ {
  next
}

!statement && line ~ /(^|[^A-Za-z0-9_$.])await[ \t(]/ && line !~ /=>|function/ && previous !~ /[=(\[,?:&|+.]$/ {
  if (previous != "" && previous !~ /[{]$/ && !(FILENAME in reported)) {
    print FILENAME
    reported[FILENAME] = 1
  }

  statement = 1
  depth = 0
}

statement {
  depth += gsub(/[({\[]/, "&", line) - gsub(/[)}\]]/, "&", line)

  if (depth <= 0 && line !~ /[=,.]$/) {
    statement = 0
    pending = 1
  }
}

{
  previous = line
}

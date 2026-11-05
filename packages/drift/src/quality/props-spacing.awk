# Prints, once, each component file whose props destructuring is not followed
# by a blank line (.claude/rules/05-react.rule.md, rule 4).
FNR == 1 {
  pending = 0
}

{
  line = $0
  sub(/^[ \t]+/, "", line)
}

pending && line != "" && line !~ /^\}/ && !(FILENAME in reported) {
  print FILENAME
  reported[FILENAME] = 1
}

{
  pending = line ~ /\} = props;?$/
}

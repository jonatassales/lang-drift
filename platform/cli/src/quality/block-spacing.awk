# Prints, once, each file with a blank line right after the { that opens a
# block or right before the } that closes one (.claude/rules/00-general.rule.md,
# Statement spacing). Comments are not code, so one ending in { opens nothing.
FNR == 1 {
  previous = ""
  blank = 0
}

{
  line = $0
  sub(/^[ \t]+/, "", line)
}

line == "" && previous ~ /\{$/ && !(FILENAME in reported) {
  print FILENAME
  reported[FILENAME] = 1
}

line ~ /^\}/ && blank && !(FILENAME in reported) {
  print FILENAME
  reported[FILENAME] = 1
}

{
  blank = line == ""
}

line != "" && line !~ /^(\/\/|\/\*|\*)/ {
  previous = line
}

grep VERSION main.go | head -n1 | awk '{ print $4 }' | tr -d '"'

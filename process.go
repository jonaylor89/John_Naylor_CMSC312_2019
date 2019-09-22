
package main

const (
  created = iota
  running
  waiting
  blocked
  terminated
)

type struct Process {
  state int,
}

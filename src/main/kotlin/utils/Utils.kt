package utils

import java.io.OutputStream
import java.io.PrintStream

fun <T> silenceConsole(lambda: () -> T): T {
  val originalErr = System.err
  val originalOut = System.out
  val noop = PrintStream(OutputStream.nullOutputStream())

  return try {
    System.setErr(noop)
    System.setOut(noop)
    lambda()
  } finally {
    noop.close()
    System.setOut(originalOut)
    System.setErr(originalErr)
  }
}


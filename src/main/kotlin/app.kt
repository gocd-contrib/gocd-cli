import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.options.default
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.required
import com.github.ajalt.clikt.parameters.types.file
import pluginResolution.loadPlugins
import java.io.File
import kotlin.system.exitProcess

fun main(args: Array<String>) = GoCD().subcommands(ConfigRepo().subcommands(ConfigRepoCheck())).main(args)

class GoCD : CliktCommand() {
  override fun run() {
    ensureDirExists(configDir())
  }
}

class ConfigRepo : CliktCommand(name = "config-repo") {
  override fun run() = Unit
}

class ConfigRepoCheck : CliktCommand(name = "check") {
  private val plugins: File by option("-d", "--plugins-dir", help = "The path from which to load plugins. Defaults to \"\$HOME/.gocd/plugins\"").file(exists = true, readable = true).default(File(configDir(), "plugins"))
  private val pluginId: String? by option("-p", "--plugin-id", help = "The config-repo plugin ID to check syntax.").required()

  override fun run() {
    ensureDirExists(File(plugins, "external"))
    val pluginRegistry = loadPlugins(plugins)

    if (!pluginRegistry.containsKey(pluginId)) {
      echo(message = "Cannot find a plugin with id `$pluginId`; known plugins in `${plugins.absolutePath}`: [${pluginRegistry.keys.joinToString(", ")}]", err = true)
      exitProcess(1)
    }

    println(pluginRegistry[pluginId])
  }
}

private fun ensureDirExists(dir: File) {
  dir.mkdirs()

  if (!dir.exists() || !dir.isDirectory) {
    System.err.println("Failed to create $dir; be sure the parent directory exists and is writable")
    exitProcess(1)
  }
}

private fun configDir(): File = File(System.getProperty("user.home"), ".gocd")
import com.github.ajalt.clikt.core.CliktCommand
import com.github.ajalt.clikt.core.subcommands
import com.github.ajalt.clikt.parameters.arguments.argument
import com.github.ajalt.clikt.parameters.options.default
import com.github.ajalt.clikt.parameters.options.option
import com.github.ajalt.clikt.parameters.options.required
import com.github.ajalt.clikt.parameters.types.file
import com.google.gson.Gson
import com.thoughtworks.go.plugin.api.request.DefaultGoPluginApiRequest
import pluginMessages.CONFIG_REPO_MESSAGE_PARSE_FILE
import pluginMessages.CONFIG_REPO_TYPE
import pluginMessages.CONFIG_REPO_VERSION
import pluginResolution.instance
import pluginResolution.loadPluginById
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
  private val pluginsDir: File by option("-d", "--plugins-dir", help = "The path from which to load pluginsDir. Defaults to \"\$HOME/.gocd/plugins\"").file(exists = true, readable = true).default(File(configDir(), "plugins"))
  private val pluginId: String by option("-p", "--plugin-id", help = "The config-repo plugin ID to check syntax.").required()
  private val file: File by argument("file", "The file to check").file(exists = true, readable = true)

  override fun run() {
    ensureDirExists(pluginsDir)

    val pluginClass = try {
      loadPluginById(pluginsDir, pluginId)
    } catch (e: Exception) {
      echo(e.message, err = true)
      exitProcess(1)
    }

    val plugin = instance(pluginId, pluginClass, CONFIG_REPO_TYPE, CONFIG_REPO_VERSION)
    val apiResponse = plugin.handle(apiRequest(file)).responseBody()

    val parse: Result =
      try {
        gson.fromJson<Result>(apiResponse, Result::class.java)
      } catch (e: Exception) {
        echo("Error occurred while parsing plugin response: ${e.message}", err = true)
        echo("Plugin responded with: ${apiResponse}", err = true)
        exitProcess(1)
      }

    when {
      parse.errors.isNotEmpty() -> {
        echo(parse.errors.joinToString("\n"), err = true)
        exitProcess(1)
      }
      else -> echo("OK")
    }
  }
}

data class ParseError(val location: String, val message: String) {
  override fun toString(): String {
    return "[$location]: $message"
  }
}

data class Result(val errors: List<ParseError>)
data class Request(val file: String)

private fun apiRequest(file: File): DefaultGoPluginApiRequest {
  val request = DefaultGoPluginApiRequest(CONFIG_REPO_TYPE, CONFIG_REPO_VERSION, CONFIG_REPO_MESSAGE_PARSE_FILE)
  request.setRequestBody(gson.toJson(Request(file.path)))
  return request
}

private fun ensureDirExists(dir: File) {
  dir.mkdirs()

  if (!dir.exists() || !dir.isDirectory) {
    System.err.println("Failed to create $dir; be sure the parent directory exists and is writable")
    exitProcess(1)
  }
}

private val gson: Gson = Gson()
private fun configDir(): File = File(System.getProperty("user.home"), ".gocd")
package pluginMessages

import com.google.gson.Gson
import com.thoughtworks.go.plugin.api.request.DefaultGoPluginApiRequest
import java.io.File

data class ParseError(val location: String, val message: String) {
  override fun toString(): String {
    return "[$location]: $message"
  }
}

data class Response(val errors: List<ParseError>)
data class Request(val file: String)

fun apiRequest(file: File): DefaultGoPluginApiRequest {
  val request = DefaultGoPluginApiRequest(CONFIG_REPO_TYPE, CONFIG_REPO_VERSION, CONFIG_REPO_MESSAGE_PARSE_FILE)
  request.setRequestBody(gson.toJson(Request(file.path)))
  return request
}

fun parseResponse(json: String): Response = gson.fromJson<Response>(json, Response::class.java)

private val gson: Gson = Gson()

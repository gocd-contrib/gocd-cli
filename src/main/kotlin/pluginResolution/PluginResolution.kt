package pluginResolution

import com.thoughtworks.go.plugin.api.GoPlugin
import io.squark.nestedjarclassloader.NestedJarClassLoader
import utils.silenceConsole
import java.io.File
import java.util.jar.JarFile
import javax.xml.parsers.DocumentBuilderFactory

fun loadPluginById(plugins: File, pluginId: String): Class<*> {
  val jarPath = jarByPluginId(plugins, pluginId)
  val jarFile = JarFile(jarPath)

  val cl = NestedJarClassLoader(ClassLoader.getSystemClassLoader(), null)
  cl.addURLs(jarPath.toURI().toURL())

  val entries = jarFile.entries()
  while (entries.hasMoreElements()) {
    val entry = entries.nextElement()

    if (!isClassEntry(entry)) continue
    val klass = cl.loadClass(toClassName(entry))

    if (meetsGoPluginCriteria(klass) && isInstantiable(klass)) {
      return klass
    }
  }

  throw IllegalStateException("Failed to identify plugin class for pluginId `$pluginId` in plugin jar [${jarPath.absolutePath}]")
}

fun instance(pluginId: String, pluginClass: Class<*>, requiredExtension: String, extensionVersion: String): GoPlugin {
  val plugin = silenceConsole {
    pluginClass.getDeclaredConstructor().newInstance() as GoPlugin
  }

  val metadata = plugin.pluginIdentifier()
  if (requiredExtension != metadata.extension) {
    throw IllegalArgumentException("Plugin `$pluginId` is not a `$requiredExtension` plugin!")
  }

  if (!metadata.supportedExtensionVersions.contains(extensionVersion)) {
    throw IllegalArgumentException("Plugin `$pluginId` must support at least $requiredExtension $extensionVersion; this one supports only ${metadata.supportedExtensionVersions.joinToString(", ")}")
  }

  return plugin
}

internal fun getPluginId(jarFile: JarFile): String? {
  val descriptor = jarFile.getEntry("plugin.xml") ?: return null
  val content = jarFile.getInputStream(descriptor)
  val xml = DocumentBuilderFactory.newInstance().newDocumentBuilder()
  val doc = xml.parse(content)
  return doc.documentElement.getAttribute("id")
}

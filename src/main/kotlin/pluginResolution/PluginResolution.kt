package pluginResolution

import com.thoughtworks.go.plugin.api.GoPlugin
import com.thoughtworks.go.plugin.api.annotation.Extension
import io.squark.nestedjarclassloader.NestedJarClassLoader
import java.io.File
import java.lang.reflect.Modifier
import java.nio.file.FileSystems
import java.nio.file.Files
import java.nio.file.Path
import java.util.jar.JarEntry
import java.util.jar.JarFile
import javax.xml.parsers.DocumentBuilderFactory
import kotlin.streams.toList

fun loadPlugins(plugins: File): Map<String, Class<*>> {
  val cl = NestedJarClassLoader(ClassLoader.getSystemClassLoader(), null)

  return Files.walk(plugins.toPath(), 1).filter(isJar()).toList().fold(mutableMapOf<String, Class<*>>()) { memo, jarPath ->
    val jarFile = JarFile(jarPath.toFile())
    cl.addURLs(jarPath.toUri().toURL())
    val id = getPluginId(jarFile)

    if (null != id && id.isNotBlank()) {
      memo[id] = GoPlugin::class.java
      val entries = jarFile.entries()
      while (entries.hasMoreElements()) {
        val entry = entries.nextElement()

        if (!isClassEntry(entry)) continue
        val klass = cl.loadClass(toClassName(entry))
        if (meetsGoPluginCriteria(klass) && isInstantiable(klass)) {
          memo[id] = klass
        }
      }
    }
    memo
  }.toMap()
}

internal fun isInstantiable(klass: Class<*>): Boolean {
  return !isANonStaticInnerClass(klass) && null != klass.getConstructor()
}

internal fun isANonStaticInnerClass(candidateClass: Class<*>): Boolean {
  return candidateClass.isMemberClass && !Modifier.isStatic(candidateClass.modifiers)
}

internal fun meetsGoPluginCriteria(klass: Class<*>): Boolean {
  return GoPlugin::class.java.isAssignableFrom(klass) &&
    null != klass.getAnnotation(Extension::class.java) &&
    !klass.isInterface && Modifier.isPublic(klass.modifiers) &&
    !Modifier.isAbstract(klass.modifiers)
}

internal fun toClassName(entry: JarEntry) =
  entry.realName.removePrefix("/").removeSuffix(".class").replace("/", ".")

internal fun isClassEntry(entry: JarEntry): Boolean {
  val fullPath = entry.realName
  return fullPath.endsWith(".class") && !(fullPath.startsWith("META-INF/") || fullPath.startsWith("lib/"))
}

internal fun getPluginId(jarFile: JarFile): String? {
  val descriptor = jarFile.getEntry("plugin.xml") ?: return null
  val content = jarFile.getInputStream(descriptor)
  val xml = DocumentBuilderFactory.newInstance().newDocumentBuilder()
  val doc = xml.parse(content)
  return doc.documentElement.getAttribute("id")
}

internal fun isJar(): (Path) -> Boolean {
  val glob = FileSystems.getDefault().getPathMatcher("glob:*.jar")
  return { p -> glob.matches(p.fileName) }
}
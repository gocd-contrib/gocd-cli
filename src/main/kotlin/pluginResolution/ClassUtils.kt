package pluginResolution

import com.thoughtworks.go.plugin.api.GoPlugin
import com.thoughtworks.go.plugin.api.annotation.Extension
import java.io.File
import java.lang.reflect.Modifier
import java.nio.file.FileSystems
import java.nio.file.Files
import java.nio.file.Path
import java.util.jar.JarEntry
import java.util.jar.JarFile


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

internal fun jarByPluginId(pluginsDir: File, pluginId: String) =
  Files.walk(pluginsDir.toPath(), 1).filter(isJar()).filter(pluginIdMatches(pluginId)).findFirst().orElseThrow { IllegalArgumentException("Failed to locate a plugin with id `$pluginId` in [${pluginsDir.absolutePath}]") }.toFile()

internal fun pluginIdMatches(pluginId: String) =
  { jarPath: Path -> pluginId.isNotBlank() && getPluginId(JarFile(jarPath.toFile())) == pluginId }

internal fun isJar(): (Path) -> Boolean {
  val glob = FileSystems.getDefault().getPathMatcher("glob:*.jar")
  return { p -> glob.matches(p.fileName) }
}
PLATFORMS = [
  [ 'darwin','amd64' ],
  [ 'freebsd', '386' ],
  [ 'freebsd', 'amd64' ],
  [ 'linux', '386' ],
  [ 'linux', 'amd64' ],
  [ 'linux', 'arm' ],
  [ 'linux', 'arm64' ],
  [ 'windows', '386' ],
  [ 'windows', 'amd64' ]
]

namespace :build do
  task :setup do
    sh 'go get -d ./...'
  end

  task :clean do
    sh 'rm -f gocd'
    sh 'rm -rf dist'
  end

  task :local, [:release] => [:clean, :setup, :test] do |t, args|
    release = "#{args[:release] || 'Edge'}-#{ENV['GO_PIPELINE_LABEL'] || "localbuild"}"
    ldflags = "-X main.Version=#{release} -X main.GitCommit=#{git_revision} -X main.Platform=$(go env GOARCH)-$(go env GOOS)"
    sh "go build -o gocd -ldflags \"#{ldflags}\""
  end

  task :prod, [:release] => [:clean, :setup, :test] do |t, args|
    release = "#{args[:release] || 'Edge'}-#{ENV['GO_PIPELINE_LABEL'] || "localbuild"}"
    PLATFORMS.each do |os, arch|
      name = "gocd"
      dir = "dist/#{os}/#{arch}"
      sh "mkdir -p #{dir}"
      if /windows/ =~ os
        name += ".exe"
        puts name
      end
      ldflags = "-X main.Version=#{release} -X main.GitCommit=#{git_revision} -X main.Platform=#{arch}-#{os}"

      sh "GOOS=#{os} GOARCH=#{arch} go build -o #{dir}/#{name} -ldflags '#{ldflags}'"
    end
  end
end

def git_revision
  `git rev-list --abbrev-commit -1 HEAD`
end

namespace :osx do
  signing_base = "codesigning"
  signing_dir = "src/osx"

  task :setup do
    cd '..' do
      mkdir_p "#{signing_base}/#{signing_dir}"
      cd 'dist/darwin/amd64' do
        sh 'chmod a+rx gocd'
        sh "zip ../../../#{signing_base}/#{signing_dir}/osx-cli.zip gocd"
      end
    end
  end

  task :sign => :setup do
    cd "../#{signing_base}" do
      sh 'rake --trace osx:sign'
      mv 'out/osx/osx-cli.zip', '..'
    end
  end
end

namespace :win do
  signing_base = "codesigning"
  signing_dir = "src/win"

  task :setup do
    cd '..' do
      mkdir_p "#{signing_base}/#{signing_dir}"
      mv 'dist/windows/amd64/gocd.exe', "#{signing_base}/#{signing_dir}"
    end
  end

  task :sign => :setup do
    cd "../#{signing_base}" do
      sh "bundle exec rake --trace win:sign"
      cd "out/win" do
        sh 'jar -cMf ../../../win-cli.zip gocd.exe'
      end
    end
  end
end

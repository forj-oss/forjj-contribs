{{/* Template is defined by ProjectModel struct (jobs_repo.go) */}}
multibranchPipelineJob('{{ .Project.Name }}') {
  description('Folder for Project {{ .Project.Name }} generated and maintained by Forjj. To update it use forjj update')
  branchSources {
{{ if eq .Project.SourceType "github" }}\
      github {
{{   if not (eq .Project.Github.ApiUrl "https://api.github.com/") }}\
          apiUri('{{ .Project.Github.ApiUrl }}')
{{   end }}\
          repoOwner('{{ .Project.Github.RepoOwner }}')
{{   if .Source.GithubUser.Name }}\
          scanCredentialsId('github-user')
{{   end }}\
          repository('{{ .Project.Github.Repo }}')
      }
{{ end }}\
{{ if eq .Project.SourceType "git" }}\
      git {
          remote('{{ .Project.Git.RemoteUrl }}')
          includes('*')
      }
{{ end }}\
  }
{{ if .Project.InfraRepo }}\
  configure {
      it / factory {
          scriptPath('apps/ci/jenkins/Jenkinsfile')
    }
  }
{{ end }}
  orphanedItemStrategy {
      discardOldItems {
          numToKeep(20)
      }
  }
}

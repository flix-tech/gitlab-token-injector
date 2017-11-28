package main

import (
    gitlab "github.com/xanzy/go-gitlab"
)

var glc *gitlab.Client

type Project struct {
    name string
    namespace string
}

func getProjects() (map[int]*Project, error) {
    var err error
    var response *gitlab.Response
    var gitlabProjects []*gitlab.Project
    var newProjects []*gitlab.Project
    projects := map[int]*Project{}
    listOpts := &gitlab.ListOptions{
        Page:       1,
        PerPage:    20,
    }
    opt := &gitlab.ListProjectsOptions{
        ListOptions:    *listOpts,
        Sort:           gitlab.String("desc"),
        OrderBy:        gitlab.String("updated_at"),
        Membership:     gitlab.Bool(true),
    }
    for opt.ListOptions.Page > 0 {
        newProjects, response, err = glc.Projects.ListProjects(opt)
        for _, tmp := range newProjects {
            gitlabProjects = append(gitlabProjects, tmp)
        }
        if err != nil {
            log.Error(err.Error())
            log.Error("Failed to get projects, please check your GitLab hostname.")
            return projects, err
        }
        log.Debugf("-- Got page: %v of %v", opt.ListOptions.Page, response.LastPage)
        opt.ListOptions.Page = response.NextPage
    }

    log.Debugf("Got %v GitLab projects...", len(gitlabProjects))

    if len(gitlabProjects) == 0 {
        log.Notice("No GitLab projects found, please check your API token and project memberships.")
        return projects, nil
    }
   
    for _, project := range gitlabProjects {
        log.Debug("-- Got project: '%v/%v' (%v)", project.Namespace.Path, project.Path, project.ID)
        tmp := new(Project)
        tmp.name = project.Path
        tmp.namespace = project.Namespace.Path
        projects[project.ID] = tmp;
    }
    return projects, nil
}


func setProjectVar(project_namespace string, project_name string, var_key string, var_value string) (bool, error) {

    log.Debugf("-- Checking if variable '%v' exists...", var_key)
    buildVariable, _, err := glc.BuildVariables.GetBuildVariable(project_namespace+"/"+project_name, var_key)

    if buildVariable != nil {
        log.Debug("-- Variable found, will update...")
        UpdateBuildVariableOptions := &gitlab.UpdateBuildVariableOptions{
            Key: gitlab.String(var_key),
            Value: gitlab.String(var_value),
            Protected: gitlab.Bool(false),
        }
        buildVariable, _, err = glc.BuildVariables.UpdateBuildVariable(project_namespace+"/"+project_name, var_key, UpdateBuildVariableOptions)
        if err == nil {
            log.Debug("-- Variable updated.")
        }
    } else {
        log.Debug("-- Variable not found, will create...")
        CreateBuildVariableOptions := &gitlab.CreateBuildVariableOptions{
            Key: gitlab.String(var_key),
            Value: gitlab.String(var_value),
            Protected: gitlab.Bool(false),
        }
        buildVariable, _, err = glc.BuildVariables.CreateBuildVariable(project_namespace+"/"+project_name, CreateBuildVariableOptions)
        if err == nil {
            log.Debug("-- Variable created.")
        }
    }
    if err != nil {
        return false, err
    }
    return true, nil
}


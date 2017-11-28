package main

import (
    "time"
    "fmt"
    "os"
    "io/ioutil"
    jwt "github.com/dgrijalva/jwt-go"
    gitlab "github.com/xanzy/go-gitlab"
    "github.com/op/go-logging"
    "github.com/namsral/flag"
)

var log = logging.MustGetLogger("default")
var logFormatDefault = logging.MustStringFormatter(
    `%{color}[%{time:2006-01-02 15:04:05}] %{color:reset} %{message}`,
)
var logFormatDebug = logging.MustStringFormatter(
    `%{color}%{time:2006-01-02T15:04:05-07:00} %{shortfunc:10.10s} ▶ %{level:6s} %{color:reset} %{message}`,
)

func main() {

    error := false
    token_ttl := flag.Duration("ttl", 48 * 3600 * time.Second, "Token lifetime in hours")
    key_path := flag.String("key", "", "Path to a private key file")
    gitlab_host := flag.String("host", "", "GitLab instance hostname")
    gitlab_token := flag.String("token", "", "GitLab API token")
    gitlab_var := flag.String("var", "KUBE_TOKEN", "GitLab secret variable")
    debug := flag.Bool("debug", false, "Enable debug output")
    flag.String(flag.DefaultConfigFlagname, "", "Path to a config file")
    flag.Parse()

    flag.Usage = func() {
        fmt.Print("\nUsage:\n\n")
        flag.PrintDefaults()
        fmt.Print("\n")
    }


    // Set up logging

    logLevel := logging.NOTICE
    logFormat := logFormatDefault

    if *debug == true {
        logLevel = logging.DEBUG
        logFormat = logFormatDebug
    }

    logBackend := logging.NewLogBackend(os.Stderr, "", 0)
    logBackendFormatter := logging.NewBackendFormatter(logBackend, logFormat)
    logBackendLeveled := logging.AddModuleLevel(logBackendFormatter)
    logBackendLeveled.SetLevel(logLevel, "")
    logging.SetBackend(logBackendLeveled)


    // Argument error handling

    if len(*key_path) == 0 {
        error = true
        log.Errorf("%s missing\n", flag.Lookup("key").Usage)
    } else {
        log.Infof("Using key path: %v", *key_path)
    }

    if len(*gitlab_host) == 0 {
        error = true
        log.Errorf("%s missing\n", flag.Lookup("host").Usage)
    } else {
        log.Infof("Using GitLab hostname: %v", *gitlab_host)
    }

    if len(*gitlab_token) == 0 {
        error = true
        log.Errorf("%s missing\n", flag.Lookup("token").Usage)
    } else {
        log.Infof("Using GitLab API token: %v", *gitlab_token)
    }

    log.Infof("Using variable name: %v", *gitlab_var)
    log.Infof("Using token TTL: %v", *token_ttl)

    if len(*key_path) > 0 {
        if _, err := os.Stat(*key_path)
        os.IsNotExist(err) {
            error = true
            log.Errorf("Private Key file '%v' does not exist.\n", *key_path)
            os.Exit(1)
        }
    }

    if error == true {
        flag.Usage()
        os.Exit(1)
    }


    // Start actually doing stuff

    privateKeyBytes, err := ioutil.ReadFile(*key_path)
    if err != nil {
        log.Fatal(err)
    }

    privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
    if err != nil {
        log.Fatal(err)
    }

    glc = gitlab.NewClient(nil, *gitlab_token)
    glc.SetBaseURL("https://" + *gitlab_host + "/api/v4")

    log.Debug("Getting GitLab projects...")
    projects, err := getProjects()

    for _, project := range projects {
        token_user := "gitlab/" + project.namespace + "/" + project.name
        token_groups := []string{"gitlab", "gitlab/" + project.namespace}

        log.Debug("Generating token for user '%v' in group '%v' with TTL %v", token_user, token_groups, token_ttl)
        jwt_token, err := generateToken(time.Now(), *token_ttl, token_user, token_groups, privateKey)

        log.Debugf("Setting variable for '%v'...", token_groups)
        success, err := setProjectVar(project.namespace, project.name, *gitlab_var, jwt_token)

        if success == false && err != nil {
            if *debug == true {
                log.Error(err.Error())
            }
            log.Errorf("%v/%v - failed -- Please make sure the GitLab user has the role 'Master' for this project.",  project.namespace, project.name)
        } else {
            log.Noticef("%v/%v - success",  project.namespace, project.name)
        }
    }

    if *debug == true {
        log.Info("Finished setting variables.")
    }

}

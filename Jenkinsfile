@Library('libpipelines@master') _

hose {
    EMAIL = 'qa'
    MODULE = 'valkiria'
    REPOSITORY = 'valkiria'
    SLACKTEAM = 'stratiopaas'
    BUILDTOOL = 'make'
    DEVTIMEOUT = 10
    LANG = 'go'
    AGENT = 'DCOS'

    DEV = { config ->        
        doCompile(config)
        doUT(config)
        doPackage(conf: config, skipOnPR: true)
        doDeploy(config)
        doStaticAnalysis(config)
     }
}

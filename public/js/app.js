var infoSys = angular.module('infoSys', ['ngRoute', 'ngAnimate', 'infoSysServices', 'infoSysDirectives', 'ngDialog']);

infoSys.run(function($rootScope) {
    angular.element(document).on("click", function(e) {
        $rootScope.$broadcast("documentClicked", angular.element(e.target));
    });
});

infoSys.config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.
            /*when('/objects', {
                templateUrl: 'public/templates/objects.html',
                controller: 'objectsCntl'
            }).
            when('/objects/:objectId', {
                templateUrl: 'public/templates/object.html',
                controller: 'objectCntl'
            }).
            when('/objects-edit', {
                templateUrl: 'public/templates/objects_edit.html',
                controller: 'entitiesCntl'
            }).
            when('/objects-edit/:objectId', {
                templateUrl: 'public/templates/object_edit.html',
                controller: 'entityCntl'
            }).
            when('/configs', {
                templateUrl: 'public/templates/configs.html',
                controller: 'configsCntl'
            }).
            when('/accounts', {
                templateUrl: 'public/templates/accounts.html',
                controller: 'accountsCntl'
            }).
            when('/groups', {
                templateUrl: 'public/templates/groups.html',
                controller: 'groupsCntl'
            }).*/
            when('/db', {
                templateUrl: 'public/templates/db.html',
                controller: 'dbCntl'
            }).
            otherwise({
                redirectTo: '/db'
            });
    }]);
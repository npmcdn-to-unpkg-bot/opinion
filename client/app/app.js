'use strict';
// Declare app level module which depends on views, and components
Storage.prototype.setObject = function(key, value) {
    this.setItem(key, JSON.stringify(value));
}
Storage.prototype.getObject = function(key) {
    var value = this.getItem(key);
    return value && JSON.parse(value);
}
window.isAdmin = function() {

    var user = localStorage.getItem("user");
    if (user) {
        return JSON.parse(user).Data.Admin;
    }
};
window.issuper = function() {

    var user = localStorage.getItem("user");



    if (user) {
        return JSON.parse(user).Data.Email == 'thesyncim@gmail.com' || JSON.parse(user).Data.Email == 'vitectv@gmail.com';
    }
};
window.isloggedin = function() {

    var user = localStorage.getItem("user");
    if (user) {
        return true;
    }
    return false;
};

window.hostname = 'http://opinion.azorestv.com/api/';
/*window.hostname = 'http://localhost:9999/';*/
window.hostnametpl = '';
angular.module('myApp', ['ngRoute', 'myApp.Articles', 'myApp.Publishers','myApp.Clients',
    'myApp.Playlist', 'myApp.Auth', 'wysiwyg.module', 'angularMoment',
    'naif.base64', 'toastr', 'angular-datepicker', 'ngCookies','mgcrea.ngStrap'
]).config(['$routeProvider',
    function($routeProvider) {
        $routeProvider.otherwise({
            redirectTo: '/articles/list'
        });
    }
]).directive('ngReallyClick', [

    function() {
        return {
            restrict: 'A',
            link: function(scope, element, attrs) {
                element.bind('click', function() {
                    var message = attrs.ngReallyMessage;
                    if (message && confirm(message)) {
                        scope.$apply(attrs.ngReallyClick);
                    }
                });
            }
        }
    }
]).config(function($httpProvider) {
    $httpProvider.defaults.withCredentials = true;
    //rest of route code
}).factory('authHttpResponseInterceptor', ['$q', '$location',
    function($q, $location) {
        return {
            response: function(response) {
                if (!localStorage.getItem("user")) {
                   $location.path('/auth/login').search('returnTo',
                        $location.path());
                }
                if (response.status === 401) {
                    console.log("Response 401");
                    localStorage.removeItem("user")
                }
                return response || $q.when(response);
            },
            responseError: function(rejection) {
                if (rejection.status === 401) {
                    console.log("Response Error 401", rejection);
                    localStorage.removeItem("user")
                }
                return $q.reject(rejection);
            }
        }
    }
]).config(['$httpProvider',
    function($httpProvider) {
        //Http Intercpetor to check auth failures for xhr requests
        $httpProvider.interceptors.push('authHttpResponseInterceptor');
    }
]).run(function($rootScope, $cookies, $cookieStore, $location) {
    $rootScope.isadmin = window.isAdmin;
    $rootScope.issuper = window.isAdmin;
    $rootScope.logout = function() {
        angular.forEach($cookies, function(v, k) {
            $cookieStore.remove(k);
        });
        localStorage.removeItem("user")
        $location.path("/")
    }
    $rootScope.isLoggedIn = window.isloggedin
}).filter('num', function() {
    return function(input) {
        return parseInt(input, 10);
    };
}).directive('scrollIf', function() {
    return function(scope, element, attributes) {
        setTimeout(function() {
            if (scope.$eval(attributes.scrollIf)) {
                window.scrollTo(0, element[0].offsetTop -
                    50)
            }
        });
    }
});
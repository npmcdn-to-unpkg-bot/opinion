'use strict';


// Declare app level module which depends on views, and components
Storage.prototype.setObject = function (key, value) {
    this.setItem(key, JSON.stringify(value));
}

Storage.prototype.getObject = function (key) {
    var value = this.getItem(key);

    return value && JSON.parse(value);
}

window.isAdmin = function () {

    var user = localStorage.getItem("user");

    if (user) {
        return JSON.parse(user).Data.Admin;
    }

};

window.isloggedin = function () {

    var user = localStorage.getItem("user");

    if (user) {
        return true;
    }
    return false;

};


window.hostname = 'http://opinion.azorestv.com/api/';
window.hostnametpl = 'http://opinion.azorestv.com/';


angular.module('myApp', [
    'ngRoute',   'myApp.Playlist',



]).config(['$routeProvider','$sceProvider', function ($routeProvider,$sceProvider) {

    $sceProvider.enabled(false);
    $routeProvider.otherwise({redirectTo: '/playlist/list'});
}]).directive('ngReallyClick', [function () {
        return {
            restrict: 'A',
            link: function (scope, element, attrs) {
                element.bind('click', function () {
                    var message = attrs.ngReallyMessage;
                    if (message && confirm(message)) {
                        scope.$apply(attrs.ngReallyClick);
                    }
                });
            }
        }
    }])


    .config(function ($httpProvider, $sceDelegateProvider) {
        $sceDelegateProvider.resourceUrlWhitelist([ 'self','**']);

        $httpProvider.defaults.withCredentials = true;
        //rest of route code
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

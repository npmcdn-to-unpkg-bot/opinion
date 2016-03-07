'use strict';



// Declare app level module which depends on views, and components
Storage.prototype.setObject = function(key, value) {
  this.setItem(key, JSON.stringify(value));
}

Storage.prototype.getObject = function(key) {
  var value = this.getItem(key);

  return value && JSON.parse(value);
}

window.isAdmin= function () {

  var user = localStorage.getItem("user");

  if(user){
    return JSON.parse(user).Data.Admin;
  }

};

window.isloggedin= function () {

  var user = localStorage.getItem("user");

  if(user){
    return true;
  }
  return false;

};



window.hostname = 'http://opinion.azorestv.com/api/';
window.hostnametpl = '/';


angular.module('myApp', [
  'ngRoute',
  'myApp.Articles',
  'angularMoment',







]).
config(['$routeProvider', function($routeProvider) {
  $routeProvider.otherwise({redirectTo: '/articles/listf'});
}]).directive('ngReallyClick', [function() {
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
    }])


.config(function ($httpProvider) {
  $httpProvider.defaults.withCredentials = true;
  //rest of route code
})

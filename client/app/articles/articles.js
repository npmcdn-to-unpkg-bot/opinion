'use strict';

angular.module('myApp.Articles', ['ngRoute'])

.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/articles/list', {
        templateUrl: window.hostnametpl+'articles/list.html',
        controller: 'ListArticles'
      });

      $routeProvider.when('/articles/create', {
        templateUrl: window.hostnametpl+'articles/create.html',
        controller: 'CreateArticles'
      });

        $routeProvider.when('/articles/edit/:id', {
            templateUrl: window.hostnametpl+'articles/edit.html',
            controller: 'EditArticles'
        });


        $routeProvider.when('/articles/listf', {
            templateUrl: window.hostnametpl+'articles/list_frontend.html',
            controller: 'ListArticlesFrontEnd'
        });

  $routeProvider.when('/articles/view/:id', {
            templateUrl: window.hostnametpl+'articles/view.html',
            controller: 'ViewArticles'
        });


}])

.controller('ListArticles', ['$scope','$http','toastr',function($scope,$http,toastr) {
        $scope.articles=[];
      var  update=function(){
            $http({
                method: 'GET',
                url:  window.hostname+'article/listall'
            }).then(function successCallback(response) {
                console.log(response)
                $scope.articles=response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update()




        $scope.delete= function(id){

            $http({
                method: 'POST',
                url:  window.hostname+'article/delete/'+id
            }).then(function successCallback(response) {
                console.log(response)
                update()
                toastr.success('Success!', 'Article Deleted');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }

        console.log($scope.articles)
}]).controller('ListArticlesFrontEnd', ['$scope','$http',function($scope,$http) {
        $scope.articles=[];
        $http({
            method: 'GET',
            url:  window.hostname+'articlef/listfrontend'
        }).then(function successCallback(response) {
            console.log(response)
            $scope.articles=response.data;
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response)
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });

        console.log($scope.articles)
    }])
.controller('CreateArticles', ['$scope','$http','toastr',function($scope,$http,toastr) {
      $scope.article={};
      $scope.article.Publisherid=  JSON.parse(window.localStorage.getItem("user")).Data.Publisherid

      $scope.submit= function(){
          $http({
              method: 'POST',
              url:  window.hostname+'article/create',
              data:$scope.article
          }).then(function successCallback(response) {
              console.log(response)
              toastr.success('Success!', 'Article Created');
              // this callback will be called asynchronously
              // when the response is available
          }, function errorCallback(response) {
              console.log(response)
              // called asynchronously if an error occurs
              // or server returns response with an error status.
          });
      }

}])

.controller('EditArticles', ['$scope','$http','$routeParams','toastr',function($scope,$http,$routeParams,toastr) {
 
        var id= $routeParams.id;
        $scope.article={};
        $scope.article.Publisherid=  JSON.parse(window.localStorage.getItem("user")).Data.Publisherid

        $http({
            method: 'GET',
            url:  window.hostname+'article/getid/'+id
            }).then(function successCallback(response) {
            console.log(response)
            $scope.article=response.data;

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response)
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });



    $scope.submit= function(){
        $http({
            method: 'POST',
            url:  window.hostname+'article/edit/'+id,
            data:$scope.article
        }).then(function successCallback(response) {
            console.log(response)
            toastr.success('Success!', 'Article Edited');
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            console.log(response)
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });
    }

}]).controller('ViewArticles', ['$scope','$http','$routeParams',function($scope,$http,$routeParams) {

           var id= $routeParams.id;
           $scope.article={};

           $http({
               method: 'GET',
               url:  window.hostname+'article/getid/'+id
           }).then(function successCallback(response) {
               console.log(response)
               $scope.article=response.data;
               // this callback will be called asynchronously
               // when the response is available
           }, function errorCallback(response) {
               console.log(response)
               // called asynchronously if an error occurs
               // or server returns response with an error status.
           });





       }]);
'use strict';

angular.module('myApp.Articles', ['ngRoute'])

    .config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/articles/list', {
            templateUrl: window.hostnametpl + 'articles/list.html',
            controller: 'ListArticles'
        });

        $routeProvider.when('/articles/create', {
            templateUrl: window.hostnametpl + 'articles/create.html',
            controller: 'CreateArticles'
        });

        $routeProvider.when('/articles/edit/:id', {
            templateUrl: window.hostnametpl + 'articles/edit.html',
            controller: 'EditArticles'
        });


        $routeProvider.when('/articles/listf', {
            templateUrl: window.hostnametpl + 'articles/list_frontend.html',
            controller: 'ListArticlesFrontEnd'
        });

        $routeProvider.when('/articles/view/:id', {
            templateUrl: window.hostnametpl + 'articles/view.html',
            controller: 'ViewArticles'
        });


    }])

    .controller('ListArticles', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.articles = [];
        var update = function () {
            $http({
                method: 'GET',
                url: window.hostname + 'article/listall'
            }).then(function successCallback(response) {

                $scope.articles = response.data;
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Failed!', response.data);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        };
        update()


        $scope.delete = function (id) {

            $http({
                method: 'POST',
                url: window.hostname + 'article/delete/' + id
            }).then(function successCallback(response) {
                console.log(response)
                update()
                toastr.success('Success!', 'Article Deleted');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Failed!', response.data);
                console.log(response)
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });

        }

        console.log($scope.articles)
    }]).controller('ListArticlesFrontEnd', ['$scope', '$http', function ($scope, $http) {
        $scope.articles = [];
        $scope.std = 'pressed';
        $scope.grid = '';
        $scope.setstd = function () {
            $scope.std = 'pressed';
            $scope.grid = '';
        }

    $scope.getPublisherImage= function (id) {

        var request = new XMLHttpRequest();
        request.open('GET', window.hostname + 'publisher/image/'+id, false);  // `false` makes the request synchronous
        request.send(null);

        if (request.status === 200) {
            console.log(request.responseText)
          return JSON.parse(request.responseText);
        }

    }

        $scope.setgrid = function () {
            $scope.std = '';
            $scope.grid = 'pressed';
        }
        $http({
            method: 'GET',
            url: window.hostname + 'articlef/listfrontend'
        }).then(function successCallback(response) {
                  $scope.articles = response.data;
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            toastr.error('Failed!', response.data);
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });

        console.log($scope.articles)
    }])
    .controller('CreateArticles', ['$scope', '$http', 'toastr', function ($scope, $http, toastr) {
        $scope.article = {};
        $scope.article.Publisherid = JSON.parse(window.localStorage.getItem("user")).Data.Publisherid;

        $scope.submit = function () {
            $http({
                method: 'POST',
                url: window.hostname + 'article/create',
                data: $scope.article
            }).then(function successCallback(response) {
                console.log(response);
                toastr.success('Success!', 'Article Created');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Failed!', response.data);
                console.log(response);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }])

    .controller('EditArticles', ['$scope', '$http', '$routeParams', 'toastr', function ($scope, $http, $routeParams, toastr) {

        var id = $routeParams.id;
        $scope.article = {};
        $scope.article.Publisherid = JSON.parse(window.localStorage.getItem("user")).Data.Publisherid;

        $http({
            method: 'GET',
            url: window.hostname + 'article/getid/' + id
        }).then(function successCallback(response) {
            console.log(response);
            $scope.article = response.data;

            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            toastr.error('Failed!', response.data);
            console.log(response);
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });


        $scope.submit = function () {
            $http({
                method: 'POST',
                url: window.hostname + 'article/edit/' + id,
                data: $scope.article
            }).then(function successCallback(response) {

                toastr.success('Success!', 'Article Edited');
                // this callback will be called asynchronously
                // when the response is available
            }, function errorCallback(response) {
                toastr.error('Failed!', response.data);
                console.log(response);
                // called asynchronously if an error occurs
                // or server returns response with an error status.
            });
        }

    }]).controller('ViewArticles', ['$scope', '$http', '$routeParams', '$sce', function ($scope, $http, $routeParams, $sce) {

    var id = $routeParams.id;
    $scope.article = {};


    var update = function () {
        $http({
            method: 'GET',
            url: window.hostname + 'article/getid/' + id
        }).then(function successCallback(response) {
            console.log(response)
            $scope.article = response.data;
            // this callback will be called asynchronously
            // when the response is available
        }, function errorCallback(response) {
            toastr.error('Failed!', response.data);
            console.log(response)
            // called asynchronously if an error occurs
            // or server returns response with an error status.
        });
    };
 
    update();

    $scope.trust = function (str) {
        return $sce.trustAsResourceUrl(str);
    };

    $scope.encode = function (str) {
        return encodeURIComponent(str);
    };

    $scope.trustTwitter = function (title) {
        return $sce.trustAsResourceUrl('https://platform.twitter.com/widgets/tweet_button.html?url=' + encodeURIComponent('http://www.azorestv.com/index.php/p/54/opiniao/#/articles/view/' + id) + '&lang=en&text=' + title)
            ;
    };

    $scope.trustFb = function () {
        var url = encodeURIComponent('http://www.azorestv.com/index.php/p/54/opiniao/#/articles/view/' + id);
        return $sce.trustAsResourceUrl('//www.facebook.com/plugins/like.php?href=' + url + '&send=false&layout=button_count&width=95&show_faces=false&action=like&colorscheme=light&font&height=21');
    };

    $scope.trustGp = function () {
        return $sce.trustAsResourceUrl('https://apis.google.com/u/0/se/0/_/+1/fastbutton?usegapi=1&size=medium&count=true&origin=http%3A%2F%2Fwww.azorestv.com&url=' + encodeURIComponent('http://www.azorestv.com/index.php/p/54/opiniao/#/articles/view/' + id) + '&gsrc=3p&jsh=m%3B%2F_%2Fscs%2Fapps-static%2F_%2Fjs%2Fk%3Doz.gapi.pt_PT.yAl1J8o04Bk.O%2Fm%3D__features__%2Fam%3DAQ%2Frt%3Dj%2Fd%3D1%2Frs%3DAGLTcCOosb3MpSoZOeJReBUqFdENTYTf3Q#_methods=onPlusOne%2C_ready%2C_close%2C_open%2C_resizeMe%2C_renderstart%2Concircled%2Cdrefresh%2Cerefresh&id=I0_1457451627556&parent=http%3A%2F%2Fwww.azorestv.com&pfname=&rpctoken=26407136')
    }
}]);
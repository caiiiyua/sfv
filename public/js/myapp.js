(function() {
	var app = angular.module('myapp', []);
	app.controller('MainController', ['$scope', '$http', function($scope, $http) {
		$scope.message = "hello";
		$scope.user_count = users.length;
		$scope.email = "";
		$scope.users = [];

		$scope.addUser = function() {
			if ($scope.email.length > 0) {
				email = $scope.email;
				$scope.users.push($scope.email);
				$scope.email = "";
				var res = $http.post('/sfvers', {'email':email});
				res.success(function(data, status, headers, config) {
					$scope.user_count = data.count;
				});
			};

		};

		$scope.getCount = function() {
			$http.get('/sfvers')
			.success(function(data) {
				$scope.user_count = data.count;
				users = data.sfvers;
				for (var i = 0; i < users.length; i++) {
					if (users[i].Email.length > 0) {
						$scope.users.push(users[i].Email)
					};
				};
			});
		};
		$scope.getCount();
	}]);

	var users = [
	];
})();
Tea.context(function () {
	this.$delay(function () {
		this.$find("form input[name='name']").focus();
	});

	this.submitSuccess = function () {
		alert("保存成功");
		window.location = "/proxy/waf/detail?wafId=" + this.config.id;
	};
});
const host = "0.0.0.0";

module.exports = {
  // transpileDependencies: ["vuetify"],
  // publicPath: "/",
  devServer: {
    host,
    hotOnly: true,
    disableHostCheck: true,
    public: "dev.homerow.tech",
    clientLogLevel: "warning",
    inline: true,
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, PATCH, OPTIONS",
      "Access-Control-Allow-Headers":
        "X-Requested-With, content-type, Authorization"
    }
  },
  configureWebpack: config => {
    config.watchOptions = {
      aggregateTimeout: 500,
      ignored: ["node_modules/**"]
      // poll: true
    };
  }
};

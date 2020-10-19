module.exports = {
  assetsDir: "ui",
  devServer: {
    proxy: {
      "/api": {
        target: "http://localhost:4433",
        ws: true,
        changeOrigin: true,
      },
    },
  },
  transpileDependencies: ["vuetify"],
};

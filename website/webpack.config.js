const path = require("path");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");
const { merge } = require("webpack-merge");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const CopyPlugin = require("copy-webpack-plugin");

module.exports = function (env) {
  const devServerConfig =
    env === "development"
      ? {
          devServer: {
            compress: true,
            port: 3002,
            host: "0.0.0.0",
            overlay: true,
            open: true,
          },
        }
      : {};

  return merge(devServerConfig, {
    mode: env,
    context: __dirname,
    entry: "./src/index.js",
    output: {
      path: path.join(__dirname, "build"),
      publicPath: "/",
      filename: "static/[name].[contenthash:8].js",
    },
    module: {
      rules: [
        {
          oneOf: [
            {
              test: /\.js$/,
              exclude: /node_modules/,
              use: {
                loader: "babel-loader",
                options: {
                  presets: ["@babel/preset-env"],
                },
              },
            },
            {
              test: /\.css$/,
              use: [
                {
                  loader: MiniCssExtractPlugin.loader,
                },
                { loader: "css-loader", options: { importLoaders: 1 } },
                "postcss-loader",
              ],
            },
            {
              test: /\.less$/,
              use: [
                {
                  loader: MiniCssExtractPlugin.loader,
                },
                { loader: "css-loader", options: { importLoaders: 3 } },
                "postcss-loader",
                "resolve-url-loader",
                "less-loader",
              ],
            },
            {
              loader: "file-loader",
              exclude: [/\.html$/, /\.js$/, /\.json$/],
              options: {
                name: "media/[name].[hash:8].[ext]",
              },
            },
          ],
        },
      ],
    },
    plugins: [
      new MiniCssExtractPlugin({
        filename: "static/[name].[contenthash:8].css",
      }),
      new CleanWebpackPlugin({
        cleanStaleWebpackAssets: false,
      }),
      new HtmlWebpackPlugin({
        template: "./public/index.html",
        cache: false,
      }),
      new CopyPlugin({
        patterns: [{ context: "./public/", from: "**/*" }],
      }),
    ],
  });
};

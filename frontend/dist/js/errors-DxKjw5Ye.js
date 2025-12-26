function hasTokenInfo(error) {
  if (typeof error !== "object" || error === null || !("tokenInfo" in error)) {
    return false;
  }
  const tokenInfo = error.tokenInfo;
  return typeof tokenInfo === "object" && tokenInfo !== null && "actual" in tokenInfo && "limit" in tokenInfo;
}
export {
  hasTokenInfo
};

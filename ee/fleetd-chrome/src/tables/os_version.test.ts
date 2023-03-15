import VirtualDatabase from "../db";

describe("os_version", () => {
  test("simple query", async () => {
    // @ts-expect-error Typescript doesn't include the userAgentData API yet.
    global.navigator.userAgentData = {
      getHighEntropyValues: jest.fn(() =>
        Promise.resolve({
          architecture: "x86",
          fullVersionList: [
            { brand: "Chromium", version: "110.0.5481.177" },
            { brand: "Not A(Brand", version: "24.0.0.0" },
            { brand: "Google Chrome", version: "110.0.5481.177" },
          ],
          mobile: false,
          model: "",
          platform: "Chrome OS",
          platformVersion: "13.2.1",
        })
      ),
    };
    chrome.runtime.getPlatformInfo = jest.fn(() =>
      Promise.resolve({ os: "cros", arch: "x86-64", nacl_arch: "x86-64" })
    );

    const db = await VirtualDatabase.init();
    const res = await db.query("select * from os_version");
    expect(res).toEqual([
      {
        name: "Chrome OS",
        platform: "chrome",
        platform_like: "chrome",
        version: "110.0.5481.177",
        major: "110",
        minor: "0",
        build: "5481",
        patch: "177",
        arch: "x86-64",
        codename: "Chrome OS 13.2.1",
      },
    ]);
  });
});

import Table from "./Table";

export default class TableGeolocation extends Table {
  name = "geolocation";
  columns = ["ip", "city", "country", "region"];

  ensureString(val: unknown): string {
    if (typeof val !== "string") {
      return val.toString();
    }
    return val;
  }

  async generate() {
    const resp = await fetch("https://ipapi.co/json");
    const json = await resp.json();
    return [
      {
        ip: this.ensureString(json.ip),
        city: this.ensureString(json.city),
        country: this.ensureString(json.country),
        region: this.ensureString(json.region),
      },
    ];
  }
}

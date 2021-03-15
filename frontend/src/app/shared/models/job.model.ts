export class Job {
    public id: string;
    public client_name: string
    public client_phone: string;
    public car_info: string;
    public appointment_info: string;
    public notes: string;
    public tag: string;
    public start: Date;
    public end: Date;
    public title: string;
    public color: any;
  }

export class Color {
  public primary: string;
  public secondary: string;
}
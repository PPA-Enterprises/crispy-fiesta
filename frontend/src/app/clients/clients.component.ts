import { Component, OnInit, ViewEncapsulation, ViewChild } from '@angular/core';
import {
  ColumnMode,
  DatatableComponent,
  SelectionType
} from '@swimlane/ngx-datatable';
import { ClientService } from '../shared/services/client.service'


export class Job {
  public id: number;
  public client_name: string
  public client_phone: string;
  public car_info: string;
  public appointment_info: string;
  public notes: string;
  public tag: string;
  public date: string;
}

export class Client {
  public id: number;
  public name: string
  public phone: string;
  public email: string;
  public jobs: Job[];
}

@Component({
  selector: 'app-page',
  templateUrl: './clients.component.html',
  styleUrls: ['./clients.component.scss', '../../assets/sass/libs/datatables.scss'],
  encapsulation: ViewEncapsulation.None
})

export class ClientsComponent {
  @ViewChild(DatatableComponent) table: DatatableComponent;
    DatatableData: Client[] = this.clientService.getAllClients();
    // row data
    public rows = this.DatatableData;
    SelectionType = SelectionType;
    selected = [];
  
    // column header
    public columns = [
      { name: 'Name', prop: 'name' },
      { name: 'Email', prop: 'email' },
      { name: 'Phone', prop: 'phone' },
    ];

    public ColumnMode = ColumnMode;

    @ViewChild('tableRowDetails') tableRowDetails: any;
    @ViewChild('tableResponsive') tableResponsive: any;
  
    // private
    private tempData = [];

    
  constructor(private clientService: ClientService) {
    this.tempData = this.DatatableData;
  }

  filterUpdate(event) {
    const val = event.target.value.toLowerCase();

    // filter our data
    const temp = this.tempData.filter(function (d) {
      return d.name.toLowerCase().indexOf(val) !== -1 || !val;
    });

    // update the rows
    this.rows = temp;
    // Whenever the filter changes, always go back to the first page
    this.table.offset = 0;
  }

  onSelect({ selected }) {
    console.log('Select Event', selected);
  }


}

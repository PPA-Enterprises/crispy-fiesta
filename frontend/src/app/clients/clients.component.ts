import { Component, OnInit, ViewEncapsulation, ViewChild } from '@angular/core';
import {
  ColumnMode,
  DatatableComponent,
  SelectionType
} from '@swimlane/ngx-datatable';
import { ClientService } from '../shared/services/client.service'
import { Job } from '../shared/models/job.model'
import { Client } from '../shared/models/client.model'


@Component({
  selector: 'app-page',
  templateUrl: './clients.component.html',
  styleUrls: ['./clients.component.scss', '../../assets/sass/libs/datatables.scss'],
  encapsulation: ViewEncapsulation.None
})

export class ClientsComponent {
  @ViewChild(DatatableComponent) table: DatatableComponent;
  DatatableData: Client[];
  private tempData = [];
  // row data
  public rows;
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

  constructor(private clientService: ClientService) {}

  ngOnInit() {
    this.clientService.getAllClients().subscribe((clients: Client[]) => {
      this.DatatableData = clients;
      this.rows = this.DatatableData;
    });
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

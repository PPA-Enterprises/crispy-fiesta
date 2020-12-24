import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { NgxDatatableModule } from '@swimlane/ngx-datatable';

import { ClientsRoutingModule } from "./clients-routing.module";

import { ClientsComponent } from "./clients.component";


@NgModule({
  imports: [
    CommonModule,
    ClientsRoutingModule,
    NgxDatatableModule
  ],
  exports: [],
  declarations: [
    ClientsComponent
  ],
  providers: [],
})
export class ClientsModule { }

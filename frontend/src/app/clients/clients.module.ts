import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";

import { ClientsRoutingModule } from "./clients-routing.module";

import { ClientsComponent } from "./clients.component";


@NgModule({
  imports: [
    CommonModule,
    ClientsRoutingModule
  ],
  exports: [],
  declarations: [
    ClientsComponent
  ],
  providers: [],
})
export class ClientsModule { }

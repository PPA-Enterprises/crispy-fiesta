import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { ClientRoutingModule } from "./client-routing.module";

import { ClientComponent } from "./client.component";


@NgModule({
  imports: [
    CommonModule,
    ClientRoutingModule,
    ReactiveFormsModule,
    FormsModule,
  ],
  exports: [],
  declarations: [
    ClientComponent,
  ],
  providers: [],
})
export class ClientModule { }

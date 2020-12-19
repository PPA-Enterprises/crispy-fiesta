import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { EditJobRoutingModule } from "./edit-job-routing.module";

import { EditJobComponent } from "./edit-job.component";


@NgModule({
  imports: [
    CommonModule,
    EditJobRoutingModule,
    ReactiveFormsModule,
    FormsModule,
  ],
  exports: [],
  declarations: [
    EditJobComponent,
  ],
  providers: [],
})
export class EditJobModule { }

import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { HomeRoutingModule } from "./home-routing.module";

import { HomeComponent } from "./home.component";


@NgModule({
  imports: [
    CommonModule,
    HomeRoutingModule,
    ReactiveFormsModule,
    FormsModule,
  ],
  exports: [],
  declarations: [
    HomeComponent,
  ],
  providers: [],
})
export class HomeModule { }

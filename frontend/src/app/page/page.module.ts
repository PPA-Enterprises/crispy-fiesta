import { NgModule } from '@angular/core';
import { CommonModule } from "@angular/common";

import { PageRoutingModule } from "./page-routing.module";

import { PageComponent } from "./page.component";


@NgModule({
  imports: [
    CommonModule,
    PageRoutingModule
  ],
  exports: [],
  declarations: [
    PageComponent
  ],
  providers: [],
})
export class PageModule { }

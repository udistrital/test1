
import { Apropiacion } from './apropiacion';

export class Movimiento {
   _id: string;
  numero:	string;
  estado_movimiento:	string;
  fecha_movimiento:	string;
  numero_oficio:	int;
  fecha_oficio:	string;
  descripcion:	string;
  unidad_ejecutora:	int;
  rubro_origen:	string;
  rubro_destino:	string;
  valor:	int;
  tipo_movimiento:	string;
  apropiacion: Apropiacion[];
}
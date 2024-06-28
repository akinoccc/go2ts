interface DemoStruct2 {
  test1: string;
  test2: number;
}

interface DemoStruct4 {
  Name: string;
  Age: number;
  Test: number;
}

interface DemoStruct3 {
  Name: string;
  Age: number;
  Test: DemoStruct;
}

interface DemoStruct {
  id_1: number;
  id_2: number;
  id_3: number;
  id_4: number;
  status: number;
  name: string;
  age: number;
  address: string;
  is_married: boolean;
  children: string[];
  salary: number;
  createdAt: string;
  updatedAt: string;
  deletedAt: string[];
  createdBy: number;
  updatedBy: number;
  deletedBy: number;
  test: DemoStruct2;
}
